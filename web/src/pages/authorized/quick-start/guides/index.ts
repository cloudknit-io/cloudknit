import { BaseGuide } from './BaseGuide';
import { ClientAccessGuide } from './ClientAccess';
import { ConfigureAWSCreds } from './ConfigureAWSCreds';
import { ConfiguringZlifecycle } from './ConfigureZlifecycle';
import { IGuide } from './IGuide';
import { OAuthGuide } from './OAuthGuide';
import { SetupEnvironmentYaml } from './SetupEnvironmentYaml';
import { SetupTeamYaml } from './SetupTeamYaml';

export const guideIndex = new Map<string, IGuide>([
    // [SetupEnvironmentYaml.getInstance().stepId, SetupEnvironmentYaml.getInstance()],
    [ConfiguringZlifecycle.getInstance().stepId, ConfiguringZlifecycle.getInstance()],
	// [ClientAccessGuide.getInstance().stepId, ClientAccessGuide.getInstance()],
    [ConfigureAWSCreds.getInstance().stepId, ConfigureAWSCreds.getInstance()],
    [SetupTeamYaml.getInstance().stepId, SetupTeamYaml.getInstance()],
    // [SetupEnvironmentYaml.getInstance().stepId, SetupEnvironmentYaml.getInstance()],
]);
export const guideKeys = [...guideIndex.keys()];
export const guideValues = [...guideIndex.values()];
