import React, { useEffect, useRef } from 'react';

const Noise = (function() {
    const p = new Uint8Array(512);
    const permutation = [ 151,160,137,91,90,15,
    131,13,201,95,96,53,194,233,7,225,140,36,103,30,69,142,8,99,37,240,21,10,23,
    190, 6,148,247,120,234,75,0,26,197,62,94,252,219,203,117,35,11,32,57,177,33,
    88,237,149,56,87,174,20,125,136,171,168, 68,175,74,165,71,134,139,48,27,166,
    77,146,158,231,83,111,229,122,60,211,133,230,220,105,92,41,55,46,245,40,244,
    102,143,54, 65,25,63,161, 1,216,80,73,209,76,132,187,208, 89,18,169,200,196,
    135,130,116,188,159,86,164,100,109,198,173,186, 3,64,52,217,226,250,124,123,
    5,202,38,147,118,126,255,82,85,212,207,206,59,227,47,16,58,17,182,189,28,42,
    223,183,170,213,119,248,152, 2,44,154,163, 70,221,153,101,155,167, 43,172,9,
    129,22,39,253, 19,98,108,110,79,113,224,232,178,185, 112,104,218,246,97,228,
    251,34,242,193,238,210,144,12,191,179,162,241, 81,51,145,235,249,14,239,107,
    49,192,214, 31,181,199,106,157,184, 84,204,176,115,121,50,45,127, 4,150,254,
    138,236,205,93,222,114,67,29,24,72,243,141,128,195,78,66,215,61,156,180
    ];
    for (let i = 0; i < 256 ; i++) {
        p[i] = p[i + 256] = permutation[i];
    }

    function fade(t: number) { return t * t * t * (t * (t * 6 - 15) + 10); }
    function lerp(t: number, a: number, b: number) { return a + t * (b - a); }
    function grad(hash: number, x: number, y: number) {
        const h = hash & 15;
        const u = h < 8 ? x : y;
        const v = h < 4 ? y : h === 12 || h === 14 ? x : 0;
        return ((h & 1) === 0 ? u : -u) + ((h & 2) === 0 ? v : -v);
    }

    return {
        get2D: function(x: number, y: number) {
            let X = Math.floor(x) & 255;
            let Y = Math.floor(y) & 255;
            x -= Math.floor(x);
            y -= Math.floor(y);
            const u = fade(x);
            const v = fade(y);
            const A = p[X] + Y, B = p[X + 1] + Y;
            return lerp(v, lerp(u, grad(p[A], x, y), grad(p[B], x - 1, y)),
                           lerp(u, grad(p[A + 1], x, y - 1), grad(p[B + 1], x - 1, y - 1)));
        }
    };
})();

const pixelScale = 6; 
const baseShimmerSpeed = 0.0007; 
const noiseSpeed = 0.0005;       
const noiseScaleX = 0.006;       
const noiseScaleY = 0.030;       

const bayerMatrix8x8 = [
    [  0, 32,  8, 40,  2, 34, 10, 42 ],
    [ 48, 16, 56, 24, 50, 18, 58, 26 ],
    [ 12, 44,  4, 36, 14, 46,  6, 38 ],
    [ 60, 28, 52, 20, 62, 30, 54, 22 ],
    [  3, 35, 11, 43,  1, 33,  9, 41 ],
    [ 51, 19, 59, 27, 49, 17, 57, 25 ],
    [ 15, 47,  7, 39, 13, 45,  5, 37 ],
    [ 63, 31, 55, 23, 61, 29, 53, 21 ]
].map(row => row.map(v => (v + 0.5) / 64));

function smoothstep(edge0: number, edge1: number, x: number) {
    const t = Math.max(0, Math.min(1, (x - edge0) / (edge1 - edge0)));
    return t * t * (3 - 2 * t);
}

export default function InputDitherBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d', { alpha: false, willReadFrequently: true });
    if (!ctx) return;

    let cw = 0;
    let ch = 0;
    let imgData: ImageData | null = null;
    let data: Uint8ClampedArray | null = null;
    let animationId: number | null = null;
    let isVisible = true;
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

    function resize() {
      if (!canvas) return;
      const rect = canvas.parentElement?.getBoundingClientRect() || canvas.getBoundingClientRect();
      const rawWidth = rect.width;
      const rawHeight = rect.height;

      cw = Math.max(1, Math.ceil(rawWidth / pixelScale));
      ch = Math.max(1, Math.ceil(rawHeight / pixelScale));
      
      canvas.width = cw;
      canvas.height = ch;
      
      if (!ctx) return;
      imgData = ctx.createImageData(cw, ch);
      data = imgData.data;

      for (let i = 0; i < data.length; i += 4) {
          data[i + 3] = 255; 
      }

      if (prefersReducedMotion.matches) {
          renderFrame(0);
      }
    }

    function renderFrame(time: number) {
      if (!imgData || !data || !ctx) return;

      const t = time;
      const shimmerSweep = t * baseShimmerSpeed;
      const noiseTimeOffsetX = t * noiseSpeed;

      for (let y = 0; y < ch; y++) {
          const matrixY = y % 8;
          const bayerRow = bayerMatrix8x8[matrixY];
          const rowStartIndex = y * cw * 4;
          const ny = y * noiseScaleY;

          for (let x = 0; x < cw; x++) {
              const nx = (x * noiseScaleX) - noiseTimeOffsetX;
              const rawNoise = Noise.get2D(nx, ny);
              const normalizedNoise = (rawNoise + 1) * 0.5;
              const waveDensity = smoothstep(0.40, 0.70, normalizedNoise);

              const shimmerWave = (Math.sin(x * 0.030 - shimmerSweep) + 1) * 0.5;
              const baseDensity = 0.05 + (shimmerWave * 0.10);
              const totalDensity = baseDensity + (waveDensity * 0.65);

              const matrixX = x % 8;
              const threshold = bayerRow[matrixX]; 

              let rgb = 10; 

              if (totalDensity > threshold) {
                  const activeAlpha = 0.04 + (waveDensity * 0.11);
                  rgb = 10 + Math.floor(235 * activeAlpha); 
              }

              const pxIdx = rowStartIndex + (x * 4);
              data[pxIdx] = rgb;     
              data[pxIdx + 1] = rgb; 
              data[pxIdx + 2] = rgb; 
          }
      }
      
      ctx.putImageData(imgData, 0, 0);
    }

    function animate(time: number) {
        if (!isVisible || prefersReducedMotion.matches) return;
        renderFrame(time);
        animationId = requestAnimationFrame(animate);
    }

    const resizeObserver = new ResizeObserver(() => resize());
    if (canvas.parentElement) {
      resizeObserver.observe(canvas.parentElement);
    } else {
      resizeObserver.observe(document.body);
    }

    function handleVisibilityChange() {
        isVisible = !document.hidden;
        if (isVisible && !prefersReducedMotion.matches) {
            if (animationId) cancelAnimationFrame(animationId);
            animationId = requestAnimationFrame(animate);
        } else if (animationId) {
            cancelAnimationFrame(animationId);
            animationId = null;
        }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange);

    const handleMotionChange = (e: MediaQueryListEvent) => {
        if (e.matches) {
            if (animationId) cancelAnimationFrame(animationId);
            renderFrame(0);
        } else {
            if (animationId) cancelAnimationFrame(animationId);
            animationId = requestAnimationFrame(animate);
        }
    };
    prefersReducedMotion.addEventListener('change', handleMotionChange);

    resize();
    if (!prefersReducedMotion.matches) {
        animationId = requestAnimationFrame(animate);
    } else {
        renderFrame(5);
    }

    return () => {
        resizeObserver.disconnect();
        document.removeEventListener('visibilitychange', handleVisibilityChange);
        prefersReducedMotion.removeEventListener('change', handleMotionChange);
        if (animationId) cancelAnimationFrame(animationId);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      className="absolute inset-0 w-full h-full z-0 pointer-events-none"
      style={{ 
        imageRendering: 'pixelated',
        maskImage: 'linear-gradient(to right, rgba(0,0,0,0.1) 0%, rgba(0,0,0,0.35) 30%, rgba(0,0,0,0.75) 65%, rgba(0,0,0,1) 100%)',
        WebkitMaskImage: 'linear-gradient(to right, rgba(0,0,0,0.1) 0%, rgba(0,0,0,0.35) 30%, rgba(0,0,0,0.75) 65%, rgba(0,0,0,1) 100%)'
      }}
      aria-hidden="true"
    />
  );
}
