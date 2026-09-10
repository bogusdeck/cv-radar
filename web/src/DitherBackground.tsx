import React, { useEffect, useRef } from 'react';

const PIXEL_SCALE = 6; 

const PALETTE = [
    [10, 10, 10],    // 0: Dark Grey
    [20, 20, 20],    // 1: Mid Grey
    [32, 32, 32],    // 2: Light Grey
    [45, 45, 45]  // 3: Highlights
];

const bayerMatrix = [
     0,  8,  2, 10,
    12,  4, 14,  6,
     3, 11,  1,  9,
    15,  7, 13,  5
].map(v => (v / 16) - 0.5);

const DITHER_SPREAD = 0.8; 

export default function DitherBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d', { alpha: false });
    if (!ctx) return;

    let time = 0;
    let animationFrameId: number | null = null;
    let canvasWidth = 0;
    let canvasHeight = 0;
    let imgData: ImageData | null = null;
    let isVisible = true;

    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

    function resizeCanvas() {
        if (!canvas) return;
        const rect = canvas.parentElement?.getBoundingClientRect() || canvas.getBoundingClientRect();
        
        canvasWidth = Math.max(1, Math.ceil(rect.width / PIXEL_SCALE));
        canvasHeight = Math.max(1, Math.ceil(rect.height / PIXEL_SCALE));
        
        canvas.width = canvasWidth;
        canvas.height = canvasHeight;
        
        if (!ctx) return;
        imgData = ctx.createImageData(canvasWidth, canvasHeight);
        
        if (prefersReducedMotion.matches) {
            renderFrame(0); 
        }
    }

    function getFieldValue(x: number, y: number, t: number) {
        let nx = x / canvasWidth;
        let ny = y / canvasHeight;

        let wave1 = Math.sin((nx * 3) + t);
        let wave2 = Math.cos((ny * 3) + (t * 0.8));
        let wave3 = Math.sin((nx + ny) * 2 - (t * 1.2));
        
        let value = (wave1 + wave2 + wave3) / 4.5 + 0.5;
        value -= ny * 0.3;

        return Math.max(0, Math.min(1, value));
    }

    function renderFrame(t: number) {
        if (!imgData || !ctx) return;
        const data = imgData.data;
        const maxColorIndex = PALETTE.length - 1;

        for (let y = 0; y < canvasHeight; y++) {
            for (let x = 0; x < canvasWidth; x++) {
                
                let val = getFieldValue(x, y, t);
                let bayerValue = bayerMatrix[(y % 4) * 4 + (x % 4)];
                let adjustedVal = val + (bayerValue * DITHER_SPREAD);

                let colorIndex = Math.round(adjustedVal * maxColorIndex);
                colorIndex = Math.max(0, Math.min(maxColorIndex, colorIndex));

                let c = PALETTE[colorIndex];
                let idx = (y * canvasWidth + x) * 4;
                data[idx] = c[0];     
                data[idx+1] = c[1];   
                data[idx+2] = c[2];   
                data[idx+3] = 255;    
            }
        }

        ctx.putImageData(imgData, 0, 0);
    }

    function animate() {
        if (!isVisible || prefersReducedMotion.matches) return;
        time += 0.03; 
        renderFrame(time);
        animationFrameId = requestAnimationFrame(animate);
    }

    const resizeObserver = new ResizeObserver(() => resizeCanvas());
    if (canvas.parentElement) {
      resizeObserver.observe(canvas.parentElement);
    } else {
      resizeObserver.observe(document.body);
    }

    function handleVisibilityChange() {
        isVisible = !document.hidden;
        if (isVisible && !prefersReducedMotion.matches) {
            animate();
        } else if (animationFrameId) {
            cancelAnimationFrame(animationFrameId);
        }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange);

    const handleMotionChange = (e: MediaQueryListEvent) => {
        if (e.matches) {
            if (animationFrameId) cancelAnimationFrame(animationFrameId);
        } else {
            animate();
        }
    };
    prefersReducedMotion.addEventListener('change', handleMotionChange);

    resizeCanvas();
    
    if (!prefersReducedMotion.matches) {
        animate();
    } else {
        renderFrame(5);
    }

    return () => {
        resizeObserver.disconnect();
        document.removeEventListener('visibilitychange', handleVisibilityChange);
        prefersReducedMotion.removeEventListener('change', handleMotionChange);
        if (animationFrameId) cancelAnimationFrame(animationFrameId);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      className="absolute inset-0 w-full h-full z-0"
      style={{ 
        imageRendering: 'pixelated',
      }}
      aria-hidden="true"
    />
  );
}
