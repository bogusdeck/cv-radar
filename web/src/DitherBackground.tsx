import React, { useEffect, useRef } from 'react';

const pixelScale = 8;
const RING_FREQ_1 = 0.08;
const RING_FREQ_2 = 0.11;
const EXPANSION_SPEED = 0.0015;

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

export default function DitherBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d', { alpha: true, willReadFrequently: true });
    if (!ctx) return;

    let cw = 0, ch = 0;
    let imgData: ImageData | null = null;
    let data: Uint8ClampedArray | null = null;
    let animationId: number | null = null;
    let isVisible = true;
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

    function resize() {
        if (!canvas) return;
        const rect = canvas.parentElement?.getBoundingClientRect() || canvas.getBoundingClientRect();
        const rawWidth = (rect.width > 0) ? rect.width : window.innerWidth;
        const rawHeight = (rect.height > 0) ? rect.height : window.innerHeight;

        cw = Math.max(1, Math.ceil(rawWidth / pixelScale));
        ch = Math.max(1, Math.ceil(rawHeight / pixelScale));
        
        canvas.width = cw;
        canvas.height = ch;
        
        if (!ctx) return;
        imgData = ctx.createImageData(cw, ch);
        data = imgData.data;

        for (let i = 0; i < data.length; i += 4) {
            data[i] = 255;     // R
            data[i + 1] = 255; // G
            data[i + 2] = 255; // B
            data[i + 3] = 0;   // A
        }

        if (prefersReducedMotion.matches) renderFrame(0);
    }

    function renderFrame(time: number) {
        if (!imgData || !data || !ctx) return;

        const t = time;
        
        const cx1 = (cw * 0.75) + Math.sin(t * 0.0004) * (cw * 0.1);
        const cy1 = (ch * 0.40) + Math.cos(t * 0.0003) * (ch * 0.1);
        
        const cx2 = (cw * 0.65) + Math.cos(t * 0.0005) * (cw * 0.15);
        const cy2 = (ch * 0.65) + Math.sin(t * 0.0006) * (ch * 0.1);

        const phaseOffset1 = t * EXPANSION_SPEED;
        const phaseOffset2 = t * EXPANSION_SPEED * 1.2;

        for (let y = 0; y < ch; y++) {
            const matrixY = y % 8;
            const bayerRow = bayerMatrix8x8[matrixY];
            const rowStartIndex = y * cw * 4;

            for (let x = 0; x < cw; x++) {
                const matrixX = x % 8;
                const threshold = bayerRow[matrixX];

                const dx1 = x - cx1;
                const dy1 = y - cy1;
                const dist1 = Math.sqrt(dx1 * dx1 + dy1 * dy1);

                const dx2 = x - cx2;
                const dy2 = y - cy2;
                const dist2 = Math.sqrt(dx2 * dx2 + dy2 * dy2);

                const wave1 = Math.sin((dist1 * RING_FREQ_1) - phaseOffset1);
                const wave2 = Math.sin((dist2 * RING_FREQ_2) - phaseOffset2);

                const normWave1 = (wave1 + 1) * 0.5;
                const normWave2 = (wave2 + 1) * 0.5;
                
                const combinedWave = (normWave1 + normWave2) * 0.5;

                const ringDensity = smoothstep(0.45, 0.75, combinedWave);

                const ambientShimmer = (Math.sin(x * 0.2 - t * 0.002) + 1) * 0.5;
                const totalDensity = 0.03 + (ambientShimmer * 0.04) + (ringDensity * 0.85);

                let alpha = 0;

                if (totalDensity > threshold) {
                    alpha = Math.floor(15 + (ringDensity * 55));
                }

                data[rowStartIndex + (x * 4) + 3] = alpha;
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
        maskImage: 'linear-gradient(to right, rgba(0,0,0,0.05) 0%, rgba(0,0,0,0.3) 35%, rgba(0,0,0,0.85) 70%, rgba(0,0,0,1) 100%)',
        WebkitMaskImage: 'linear-gradient(to right, rgba(0,0,0,0.05) 0%, rgba(0,0,0,0.3) 35%, rgba(0,0,0,0.85) 70%, rgba(0,0,0,1) 100%)'
      }}
      aria-hidden="true"
    />
  );
}
