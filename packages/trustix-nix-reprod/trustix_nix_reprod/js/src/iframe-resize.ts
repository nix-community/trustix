export function resizeIFrames(): void {
  document.querySelectorAll("iframe.seamless").forEach((elem: HTMLIFrameElement) => {
    elem.style.height = elem.contentWindow.document.body.scrollHeight + 50 + 'px';
  })
}
