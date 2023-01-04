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
  startDateTime: string;
}

export class UpdateComponentReconciliationDto {
  status?: string;
  endDateTime?: string;
  approvedBy?: string;
}
