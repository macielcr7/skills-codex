import { Module } from '@nestjs/common'

import { CommonModule } from '#/infrastructure/nest/common/common.module'
import { MediaModule } from '#/infrastructure/nest/media/media.module'

@Module({
  imports: [CommonModule, MediaModule],
})
export class AppModule {}
