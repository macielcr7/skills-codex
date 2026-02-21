import type { Video } from '#/domain/entity/media/video/video'
import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'

export class MemoryVideoRepository implements VideoRepository {
  private readonly videos = new Map<string, Video>()

  async save(video: Video): Promise<void> {
    this.videos.set(video.id, video)
  }

  async getByID(id: string): Promise<Video | null> {
    return this.videos.get(id) ?? null
  }
}
