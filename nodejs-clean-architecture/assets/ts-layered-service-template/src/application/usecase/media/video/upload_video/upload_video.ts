import { Video } from '#/domain/entity/media/video/video'
import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'

import type { IdGenerator } from '#/application/service/shared/id_generator.interface'
import type {
  UploadVideoInput,
  UploadVideoOutput,
} from '#/application/usecase/media/video/upload_video/upload_video.type'

export class UploadVideo {
  constructor(
    private readonly repository: VideoRepository,
    private readonly idGenerator: IdGenerator,
  ) {}

  async execute(input: UploadVideoInput): Promise<UploadVideoOutput> {
    const video = new Video({ id: this.idGenerator.newID(), title: input.title })
    video.validate()

    await this.repository.save(video)

    return { id: video.id }
  }
}
