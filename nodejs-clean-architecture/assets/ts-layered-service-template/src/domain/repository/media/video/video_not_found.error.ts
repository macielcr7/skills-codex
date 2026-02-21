export class VideoNotFoundError extends Error {
  constructor(id: string) {
    super(`video not found: ${id}`)
    this.name = 'VideoNotFoundError'
  }
}

