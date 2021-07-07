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

    let savedEntry = null
    if (reconcileId) {
      const existingEntry = await this.environmentReconcileRepository.findOne(
        reconcileId,
      )
      existingEntry.end_date_time = runData.endDateTime
      existingEntry.status = runData.status
      savedEntry = await this.environmentReconcileRepository.save(existingEntry)
    } else {
      const entry: EnvironmentReconcile = {
        reconcile_id: reconcileId,
        name: runData.name,
        start_date_time: runData.startDateTime,
        team_name: runData.teamName,
        status: runData.status,
        end_date_time: runData.endDateTime,
      }
      savedEntry = await this.environmentReconcileRepository.save(entry)
    }
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

    let componentEntry: ComponentReconcile = Mapper.mapToComponentReconcile(
      savedEntry,
      runData.componentReconciles,
    )[0];
    
    if (componentEntry.reconcile_id) {
      const existingEntry = await this.componentReconcileRepository.findOne(componentEntry.reconcile_id);
      existingEntry.end_date_time = componentEntry.end_date_time;
      existingEntry.status = componentEntry.status;
      componentEntry = existingEntry;
    }

    const entry = await this.componentReconcileRepository.save(componentEntry);
    return entry.reconcile_id;
  }
}
