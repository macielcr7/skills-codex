import { Module } from '@nestjs/common'

import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'

import { GetVideo } from '#/application/usecase/media/video/get_video/get_video'
import { UploadVideo } from '#/application/usecase/media/video/upload_video/upload_video'
import { GetVideoController } from '#/infrastructure/nest/media/controllers/video/get_video.controller'
import { UploadVideoController } from '#/infrastructure/nest/media/controllers/video/upload_video.controller'
import { CommonModule } from '#/infrastructure/nest/common/common.module'
import { UUIDGenerator } from '#/infrastructure/nest/common/services/uuid_generator'
import { MEDIA_TOKENS } from '#/infrastructure/nest/media/media.tokens'
import { MemoryVideoRepository } from '#/infrastructure/repository/media/video/memory/video.repository'

@Module({
  imports: [CommonModule],
  controllers: [UploadVideoController, GetVideoController],
  providers: [
    {
      provide: MEDIA_TOKENS.videoRepository,
      useClass: MemoryVideoRepository,
    },
    {
      provide: UploadVideo,
      useFactory: (repo: VideoRepository, idGen: UUIDGenerator) =>
        new UploadVideo(repo, idGen),
      inject: [MEDIA_TOKENS.videoRepository, UUIDGenerator],
    },
    {
      provide: GetVideo,
      useFactory: (repo: VideoRepository) => new GetVideo(repo),
      inject: [MEDIA_TOKENS.videoRepository],
    },
  ],
})
export class MediaModule {}

