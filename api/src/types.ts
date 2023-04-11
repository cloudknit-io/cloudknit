import { applyDecorators } from '@nestjs/common';
import { ApiParam, ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsNumber, IsString } from 'class-validator';
import { Request } from 'express';
import { Component, Environment, EnvironmentReconcile, Team } from './typeorm';
import { Organization } from './typeorm/Organization.entity';

export type APIRequest = Request & {
  org: Organization;
  team: Team;
  env?: Environment;
  argoCDAuthHeader?: string;
};

export const SqlErrorCodes = {
  NO_DEFAULT: 'ER_NO_DEFAULT_FOR_FIELD',
  DUP_ENTRY: 'ER_DUP_ENTRY',
};

export class BaseApiError extends Error {}
export class NoEnvironmentsError extends BaseApiError {}

export function OrgApiParam(): MethodDecorator {
  return applyDecorators(ApiParam({ name: 'orgId', required: true }));
}

export function TeamApiParam(): MethodDecorator {
  return applyDecorators(
    OrgApiParam(),
    ApiParam({ name: 'teamId', required: true })
  );
}

export function EnvironmentApiParam(): MethodDecorator {
  return applyDecorators(
    TeamApiParam(),
    ApiParam({ name: 'environmentId', required: true })
  );
}

export enum InternalEventType {
  ComponentCostUpdate = 'ComponentEntity.update.cost',
  EnvironmentCostUpdate = 'EnvironmentEntity.update.cost',
  EnvironmentReconCostUpdate = 'EnvironmentRecon.update.cost',
  EnvironmentReconEnvUpdate = 'EnvironmentRecon.update.env',
}

export interface InternalEvent {
  type: InternalEventType;
  payload: any;
}

export class ComponentCostUpdateEvent implements InternalEvent {
  type: InternalEventType;
  payload: Component;

  constructor(p: Component) {
    this.type = InternalEventType.ComponentCostUpdate;
    this.payload = p;
  }
}

export class EnvironmentReconCostUpdateEvent implements InternalEvent {
  type: InternalEventType;
  payload: EnvironmentReconcile;

  constructor(p: EnvironmentReconcile) {
    this.type = InternalEventType.EnvironmentReconCostUpdate;
    this.payload = p;
  }
}

export class EnvironmentReconEnvUpdateEvent implements InternalEvent {
  type: InternalEventType;
  payload: EnvironmentReconcile;

  constructor(p: EnvironmentReconcile) {
    this.type = InternalEventType.EnvironmentReconEnvUpdate;
    this.payload = p;
  }
}

export class EnvironmentCostUpdateEvent implements InternalEvent {
  type: InternalEventType;
  payload: Environment;

  constructor(p: Environment) {
    this.type = InternalEventType.EnvironmentCostUpdate;
    this.payload = p;
  }
}

export class ApiHttpException {
  @ApiProperty()
  @IsNumber()
  @IsNotEmpty()
  statusCode: number;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  message: string;
}
