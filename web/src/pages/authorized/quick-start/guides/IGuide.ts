import { NotificationsManager } from "components/argo-core";
import React from "react";

export interface IGuide {
    stepId: string;
    stepName: string;
    render: (baseClassName: string, ctx: any, nm?: NotificationsManager) => JSX.Element;
    onNext?: () => Promise<any>;
    onFinish?: () => Promise<any>;
}