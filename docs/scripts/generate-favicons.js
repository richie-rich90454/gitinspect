import sharp from 'sharp'
import { readFileSync, mkdirSync, existsSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const publicDir = join(__dirname, '..', 'public')
const svgPath = join(publicDir, 'logo.svg')

if (!existsSync(svgPath)) {
  console.error('logo.svg not found at', svgPath)
  process.exit(1)
}

const svgBuffer = readFileSync(svgPath)

const sizes = [
  { name: 'favicon-16x16.png', size: 16 },
  { name: 'favicon-32x32.png', size: 32 },
  { name: 'favicon-48x48.png', size: 48 },
  { name: 'apple-touch-icon.png', size: 180 },
  { name: 'android-chrome-192x192.png', size: 192 },
  { name: 'android-chrome-512x512.png', size: 512 },
]

async function generatePngs() {
  for (const { name, size } of sizes) {
    const outPath = join(publicDir, name)
    await sharp(svgBuffer)
      .resize(size, size)
      .png()
      .toFile(outPath)
    console.log(`  ✓ ${name} (${size}x${size})`)
  }
}

async function generateIco() {
  const png16 = await sharp(svgBuffer).resize(16, 16).png().toBuffer()
  const png32 = await sharp(svgBuffer).resize(32, 32).png().toBuffer()
  const png48 = await sharp(svgBuffer).resize(48, 48).png().toBuffer()

  const images = [
    { width: 16, height: 16, data: png16 },
    { width: 32, height: 32, data: png32 },
    { width: 48, height: 48, data: png48 },
  ]

  const numImages = images.length
  const headerSize = 6
  const dirEntrySize = 16
  const dirSize = dirEntrySize * numImages
  let dataOffset = headerSize + dirSize

  const dirEntries = []
  const imageDataChunks = []

  for (const img of images) {
    const dataSize = img.data.length
    dirEntries.push({
      width: img.width === 256 ? 0 : img.width,
      height: img.height === 256 ? 0 : img.height,
      dataSize,
      dataOffset,
    })
    imageDataChunks.push(img.data)
    dataOffset += dataSize
  }

  const totalSize = dataOffset
  const buf = Buffer.alloc(totalSize)
  let offset = 0

  buf.writeUInt16LE(0, offset); offset += 2
  buf.writeUInt16LE(1, offset); offset += 2
  buf.writeUInt16LE(numImages, offset); offset += 2

  for (let i = 0; i < numImages; i++) {
    const e = dirEntries[i]
    buf.writeUInt8(e.width, offset); offset += 1
    buf.writeUInt8(e.height, offset); offset += 1
    buf.writeUInt8(0, offset); offset += 1
    buf.writeUInt8(0, offset); offset += 1
    buf.writeUInt16LE(1, offset); offset += 2
    buf.writeUInt16LE(32, offset); offset += 2
    buf.writeUInt32LE(e.dataSize, offset); offset += 4
    buf.writeUInt32LE(e.dataOffset, offset); offset += 4
  }

  for (const chunk of imageDataChunks) {
    chunk.copy(buf, offset)
    offset += chunk.length
  }

  const icoPath = join(publicDir, 'favicon.ico')
  const { writeFileSync } = await import('fs')
  writeFileSync(icoPath, buf)
  console.log('  ✓ favicon.ico (16+32+48)')
}

console.log('Generating favicons from logo.svg...\n')
await generatePngs()
await generateIco()
console.log('\nDone! All favicons generated in docs/public/')
