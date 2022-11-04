import { configure } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { streamMapper } from 'helpers/streamMapper';
import { ApplicationStatus, ApplicationWatchEvent } from 'models/argo.models';
import { TeamItem } from 'models/projects.models';
import { ArgoMapper } from 'services/argo/ArgoMapper';

const mockProjects: TeamItem[] = [
	{
		id: 'demo',
		labels: { type: 'project' },
		healthStatus: 'Healthy',
		syncStatus: 'Synced',
		resourceVersion: '11586186',
		name: 'demo',
	},
	{
		id: 'demo-1',
		labels: { type: 'project' },
		healthStatus: 'Healthy',
		syncStatus: 'Synced',
		resourceVersion: '11586182',
		name: 'demo-1',
	},
	{
		id: 'demo-2',
		labels: { type: 'project' },
		healthStatus: 'Healthy',
		syncStatus: 'Synced',
		resourceVersion: '11586183',
		name: 'demo-2',
	},
];

const commonMockStreamStatus: ApplicationStatus = {
	status: {
		sync: {
			status: 'Unknown',
		},
		health: {
			status: 'Healthy',
		},
	},
};

const mockStream = {
	type: 'MODIFIED',
	application: {
		metadata: {
			name: 'demo',
			labels: {
				type: 'project',
			},
		},
		status: {
			sync: {
				status: 'Unknown',
			},
			health: {
				status: 'Healthy',
			},
		},
	},
};

const mockStreamAddedLabel: Partial<ApplicationWatchEvent> = {
	type: 'MODIFIED',
	application: {
		metadata: {
			name: 'demo',
			labels: {
				type: 'project',
				test: 'demo',
			},
		},
		...commonMockStreamStatus,
	},
};

const mockStreamDeletedDemoProject: Partial<ApplicationWatchEvent> | null = {
	type: 'DELETED',
	application: {
		metadata: {
			id: 'demo',
			name: 'demo',
			labels: {
				type: 'project',
			},
		},
		...commonMockStreamStatus,
	},
};

const mockStreamAddsDemo3Project: Partial<ApplicationWatchEvent> = {
	type: 'ADDED',
	application: {
		metadata: {
			name: 'demo-3',
			labels: {
				type: 'project',
			},
		},
		...commonMockStreamStatus,
	},
};

configure({ adapter: new Adapter() });
describe('Stream mapper', () => {
	it('should validate if there is stream data', () => {
		expect(streamMapper(null, mockProjects, ArgoMapper.parseTeam, 'project')).toEqual(mockProjects);
	});

	it('should return different list of projects after stream update', () => {
		expect(streamMapper(mockStream, mockProjects, ArgoMapper.parseTeam, 'project')).not.toEqual(mockProjects);
	});

	it('should add label to demo project', () => {
		const newMockProjects = streamMapper(mockStreamAddedLabel, mockProjects, ArgoMapper.parseTeam, 'project');
		expect(newMockProjects[0].labels).toEqual(mockStreamAddedLabel.application.metadata.labels);
	});

	it('should delete demo project', () => {
		const newMockProjects = streamMapper(
			mockStreamDeletedDemoProject,
			mockProjects,
			ArgoMapper.parseTeam,
			'project'
		);
		expect(newMockProjects.length).toEqual(2);
	});

	it('should add demo-3 project', () => {
		const newMockProjects = streamMapper(mockStreamAddsDemo3Project, mockProjects, ArgoMapper.parseTeam, 'project');
		expect(newMockProjects.length).toEqual(3);
	});
});
