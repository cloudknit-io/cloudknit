import { Module } from '@nestjs/common';
import { ComponentService } from './component.service';
import { ComponentController } from './component.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component } from 'src/typeorm';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component
    ])
  ],
  controllers: [ComponentController],
  providers: [ComponentService]
})
export class ComponentModule {}
