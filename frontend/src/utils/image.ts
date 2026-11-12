/** Largest edge of an uploaded image; bigger files are scaled down. */
const MAX_DIMENSION = 1200
const WEBP_QUALITY = 0.85

/**
 * Scales an uploaded image down to `maxSize` on its longest edge and encodes
 * it as a WebP data URL so it survives in localStorage. Falls back to PNG
 * when the browser cannot encode WebP.
 */
export async function createImageDataUrl(file: File, maxSize = MAX_DIMENSION): Promise<string> {
  const bitmap = await createImageBitmap(file)

  try {
    const scale = Math.min(1, maxSize / Math.max(bitmap.width, bitmap.height))
    const width = Math.max(1, Math.round(bitmap.width * scale))
    const height = Math.max(1, Math.round(bitmap.height * scale))

    const canvas = document.createElement('canvas')
    canvas.width = width
    canvas.height = height

    const context = canvas.getContext('2d')

    if (!context) throw new Error('Canvas is not supported')

    context.drawImage(bitmap, 0, 0, width, height)

    const webp = canvas.toDataURL('image/webp', WEBP_QUALITY)

    return webp.startsWith('data:image/webp') ? webp : canvas.toDataURL('image/png')
  } finally {
    bitmap.close()
  }
}
