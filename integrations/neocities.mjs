import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

export default function neocitiesIntegration() {
    return {
        name: 'neocities',
        hooks: {
            'astro:build:done': ({ dir }) => {
                const dirPath = typeof dir === 'string' ? dir : fileURLToPath(dir);
                const notFoundDirPath = path.join(dirPath, 'not_found');
                const notFoundIndexPath = path.join(notFoundDirPath, 'index.html');
                const targetPath = path.join(dirPath, 'not_found.html');

                if (fs.existsSync(notFoundIndexPath)) {
                    fs.copyFileSync(notFoundIndexPath, targetPath);
                    fs.rmSync(notFoundDirPath, { recursive: true, force: true });
                    console.log('✅ Created not_found.html for Neocities');
                } else {
                    console.log('⚠️ not_found directory not found, skipping Neocities integration');
                }
            }
        }
    };
}