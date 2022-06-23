import { ComponentReconcile } from "./component-reconcile.entity";
import { Component } from "./component.entity";
import { EnvironmentReconcile } from "./environment-reconcile.entity";
import { Environment } from "./environment.entity";

export const reconcileEntities = [
  EnvironmentReconcile,
  ComponentReconcile,
  Component,
  Environment
];
