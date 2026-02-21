import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'
import { VideoNotFoundError } from '#/domain/repository/media/video/video_not_found.error'
import type { GetVideoInput, GetVideoOutput } from '#/application/usecase/media/video/get_video/get_video.type'

export class GetVideo {
  constructor(private readonly repository: VideoRepository) {}

  async execute(input: GetVideoInput): Promise<GetVideoOutput> {
    const video = await this.repository.getByID(input.id)
    if (!video) {
      throw new VideoNotFoundError(input.id)
    }

    return { id: video.id, title: video.title }
  }
}
