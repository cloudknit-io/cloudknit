import { ConfigureAWSCreds } from './ConfigureAWSCreds';
import { ConfiguringZlifecycle } from './ConfigureZlifecycle';
import { IGuide } from './IGuide';
import { SetupTeamYaml } from './SetupTeamYaml';

export const guideIndex = new Map<string, IGuide>([
    [ConfiguringZlifecycle.getInstance().stepId, ConfiguringZlifecycle.getInstance()],
    [ConfigureAWSCreds.getInstance().stepId, ConfigureAWSCreds.getInstance()],
    [SetupTeamYaml.getInstance().stepId, SetupTeamYaml.getInstance()],
]);
export const guideKeys = [...guideIndex.keys()];
export const guideValues = [...guideIndex.values()];
