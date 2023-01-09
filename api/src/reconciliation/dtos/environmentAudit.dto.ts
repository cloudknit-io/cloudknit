import { EnvironmentReconcile } from "src/typeorm";

export interface EnvironmentReconcileWrap extends EnvironmentReconcile {
    duration: number;
}
