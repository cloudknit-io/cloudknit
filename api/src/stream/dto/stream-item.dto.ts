import { ComponentReconcileWrap } from 'src/reconciliation/dtos/componentAudit.dto';
import { EnvironmentReconcileWrap } from 'src/reconciliation/dtos/environmentAudit.dto';
import { Component, Environment, Team } from 'src/typeorm';

export class StreamItem {
  data:
    | Team
    | Environment
    | Component
    | ComponentReconcileWrap
    | EnvironmentReconcileWrap;
  type: StreamTypeEnum;
}

export enum StreamTypeEnum {
  Team = 'Team',
  Environment = 'Environment',
  Component = 'Component',
  ComponentReconcile = 'ComponentReconcile',
  EnvironmentReconcile = 'EnvironmentReconcile',
  Empty = 'Empty',
}
