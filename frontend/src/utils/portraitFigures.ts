/**
 * Vertical figure images are letterboxed over a blurred backdrop of the same
 * photo instead of being cropped to the 16:9 frame (see `article.css`). The
 * markdown renderer cannot know an image's aspect ratio in advance, so every
 * figure is marked at runtime once its main image has dimensions — either
 * right away for cached images or on the `load` event.
 *
 * Call this after the rendered HTML lands in the DOM (article page and editor
 * preview both re-run it whenever their `v-html` content changes).
 */
export function detectPortraitFigures(root: ParentNode): void {
  for (const figure of root.querySelectorAll<HTMLElement>('figure')) {
    const img = figure.querySelector<HTMLImageElement>('.figure__img')

    if (!img) continue

    const apply = () => {
      figure.classList.toggle('figure--portrait', img.naturalHeight > img.naturalWidth)
    }

    if (img.complete && img.naturalWidth > 0) {
      apply()
    } else {
      img.addEventListener('load', apply, { once: true })
    }
  }
}
