export class CreateEnvironmentReconciliationDto {
  name: string;
  teamName: string;
  startDateTime: string;
}

export class UpdateEnvironmentReconciliationDto {
  status: string;
  endDateTime?: string;
}

export class CreateComponentReconciliationDto {
  envReconcileId: number;
  name: string;
  envName: string;
  teamName: string;
  startDateTime: string;
}

export class UpdateComponentReconciliationDto extends CreateComponentReconciliationDto {
  status?: string;
  endDateTime?: string;
  approvedBy?: string;
}
