import { Injectable } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Mapper } from 'src/costing/utilities/mapper'
import { ComponentReconcile } from 'src/typeorm/reconciliation/component-reconcile.entity'
import { EnvironmentReconcile } from 'src/typeorm/reconciliation/environment-reconcile.entity'
import { Repository } from 'typeorm/repository/Repository'
import { EvnironmentReconcileDto } from '../dtos/reconcile.Dto'

@Injectable()
export class ReconciliationService {
  constructor(
    @InjectRepository(EnvironmentReconcile)
    private readonly environmentReconcileRepository: Repository<
      EnvironmentReconcile
    >,
    @InjectRepository(ComponentReconcile)
    private readonly componentReconcileRepository: Repository<
      ComponentReconcile
    >,
  ) {}

  async saveOrUpdateEnvironment(runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId)
    const entry: EnvironmentReconcile = {
      reconcile_id: reconcileId,
      name: runData.name,
      start_date_time: runData.startDateTime,
      team_name: runData.teamName,
      status: runData.status,
      end_date_time: runData.endDateTime,
    }
    const savedEntry = await this.environmentReconcileRepository.save(entry)
    return savedEntry.reconcile_id
  }

  async saveOrUpdateComponent(runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId)
    if (!reconcileId) {
      throw 'Reconcile Id is mandatory to save or update component'
    }
    const savedEntry = await this.environmentReconcileRepository.findOne(
      reconcileId,
    )
    const componentEntries: ComponentReconcile[] = Mapper.mapToComponentReconcile(
      savedEntry,
      runData.componentReconciles,
    )
    await this.componentReconcileRepository.save(componentEntries)
    return componentEntries[0].reconcile_id;
  }
}
