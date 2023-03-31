import { Component } from './component.entity';
import { User } from './User.entity';
import { Organization } from './Organization.entity';
import { ComponentReconcile } from './component-reconcile.entity';
import { EnvironmentReconcile } from './environment-reconcile.entity';
import { Environment } from './environment.entity';
import { Team } from './team.entity';
import { get } from 'src/config';

export const entities = [
  User,
  Organization,
  Component,
  ComponentReconcile,
  EnvironmentReconcile,
  Environment,
  Team,
];

const dbConfig = {
  type: 'mysql',
  host: get().TypeORM.host,
  port: get().TypeORM.port,
  username: get().TypeORM.username,
  password: get().TypeORM.password,
  database: get().TypeORM.database,
  entities,
  migrations: ['src/typeorm/migrations/*.js'],
  migrationsRun: true,
  synchronize: true,
  logging: ['error', 'schema'],
};

export {
  dbConfig,
  User,
  Organization,
  Component,
  ComponentReconcile,
  EnvironmentReconcile,
  Environment,
  Team,
};
