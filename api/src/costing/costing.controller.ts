import { Controller, Get, Query, Request } from '@nestjs/common'
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from 'src/reconciliation/validationPipes';
import { ComponentDto } from './dtos/Component.dto';
import { ComponentService } from './services/component.service'

@Controller({
  version: '1'
})
export class CostingController {
  constructor(
    private readonly compSvc: ComponentService,
  ) {}

}
