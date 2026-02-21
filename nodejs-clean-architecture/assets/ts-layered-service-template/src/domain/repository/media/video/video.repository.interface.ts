import type { Video } from '#/domain/entity/media/video/video'

export interface VideoRepository {
  save(video: Video): Promise<void>
  getByID(id: string): Promise<Video | null>
}
