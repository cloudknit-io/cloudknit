import { IsDateString, IsNotEmpty, IsNumber, IsOptional, IsString } from "class-validator";

export class CreateEnvironmentReconciliationDto {
  @IsString()
  @IsNotEmpty()
  name: string;
  
  @IsString()
  @IsNotEmpty()
  teamName: string;

  @IsDateString()
  @IsNotEmpty()
  startDateTime: string;
}

export class UpdateEnvironmentReconciliationDto {
  @IsString()
  @IsNotEmpty()
  status: string;

  @IsDateString()
  @IsOptional()
  endDateTime?: string;
}

export class CreateComponentReconciliationDto {
  @IsNumber()
  @IsNotEmpty()
  envReconcileId: number;

  @IsString()
  @IsNotEmpty()
  name: string;

  @IsDateString()
  @IsNotEmpty()
  startDateTime: string;
}

export class UpdateComponentReconciliationDto {
  @IsString()
  @IsOptional()
  status?: string;

  @IsDateString()
  @IsOptional()
  endDateTime?: string;

  @IsString()
  @IsOptional()
  approvedBy?: string;
}
