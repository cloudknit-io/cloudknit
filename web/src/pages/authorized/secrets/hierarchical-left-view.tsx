import { useApi } from 'hooks/use-api/useApi';
import React, { useEffect, useMemo, useState } from 'react';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';
import { ReactComponent as DropdownArrow } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as Add } from 'assets/images/icons/add.svg';
import { SecretsService } from 'services/secrets/secrets.service';
import { Subject } from 'rxjs';
import AuthStore from 'auth/AuthStore';
import { EntityStore } from 'models/entity.store';

export type Hierarchy = {
	id: string;
	name: string | JSX.Element;
	children: Hierarchy[];
	selectable?: boolean;
	expanded?: boolean;
	type: 'SECRET' | 'TAB';
	updateCallback?: () => void;
	_par?: Hierarchy;
};

export type Props = {
	hierarchyChanged: (h: Hierarchy) => boolean;
	refreshView: Subject<string>;
};

export const HierachicalLeftView: React.FC<Props> = ({ hierarchyChanged, refreshView }) => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const org = AuthStore.getOrganization();
	const [hierarchy, setHeirarchy] = useState<Hierarchy[]>([]);
	const [selectedHierarchy, setSelectedHierarchy] = useState<Hierarchy>();
	const [environmentsMap, setEnvironmentsMap] = useState<Map<string, string[]>>(new Map<string, string[]>());
	const selectHierarchy = (h: Hierarchy) => {
		const selected = hierarchyChanged(h);
		setSelectedHierarchy(h);
	};

	const getHierarchy = (
		id: string,
		children: Hierarchy[],
		name: string | JSX.Element,
		selectable?: boolean,
		updateCallback?: () => void,
		_par?: Hierarchy,
		type?: 'SECRET' | 'TAB'
	): Hierarchy => ({
		id,
		children,
		name,
		selectable,
		updateCallback,
		_par,
		type: type || 'SECRET',
	});

	useEffect(() => {
		if (!refreshView) return;
		const sub = refreshView.subscribe((id: string) => {
			const par = selectedHierarchy?._par;
			if (!par) return;
			id = id.replaceAll('/', ':');
			let newEnv = par.children.find(e => e.id === id) as Hierarchy;
			if (!newEnv) {
				newEnv = getHierarchy(id, [], id.replace(par?.id + ':', ''), true, undefined, par);
				par.children.push(newEnv);
				setHeirarchy([...hierarchy]);
			}

			setTimeout(() => {
				selectHierarchy(newEnv);
				const li = document.getElementById(id) as HTMLLIElement;
				li.scrollIntoView({ behavior: 'smooth' });
			}, 500);
		});

		return () => sub.unsubscribe();
	}, [selectedHierarchy]);

	useEffect(() => {
		const subscription = entityStore.emitter.subscribe(update => {
			const teams = update.teams;
			const envs = update.environments;
			if (teams.length === 0 && envs.length === 0) return;
			const map = new Map<string, string[]>();
			teams.forEach(e =>
				map.set(
					e.name,
					envs.filter(env => env.teamId === e.id).map(env => env.name)
				)
			);
			setEnvironmentsMap(map)
		});

		return () => {
			subscription.unsubscribe();
		};
	}, []);

	useEffect(() => {
		if (!org) {
			return;
		}
		const hierarchy: Hierarchy[] = [getHierarchy(org.name || '', [], 'Global Secrets', true, () => {})];

		hierarchy.push(
			getHierarchy(
				'Teams',
				[...environmentsMap.keys()].map(e => {
					const team = getHierarchy(`${org.name || ''}:${e}`, [], e, true, () => {});
					team.children = [
						getHierarchy(
							`${org.name || ''}:${e}:`,
							[],
							<span className="d-flex align-center">
								New <Add style={{ marginLeft: '5px' }} />
							</span>,
							true,
							() => {},
							team
						),
						...(environmentsMap.get(e) || [])
							.sort()
							.map(r =>
								getHierarchy(
									`${org.name || ''}:${e}:${r.replace(e + '-', '')}`,
									[],
									r.replace(e + '-', ''),
									true,
									() => {},
									team
								)
							),
					];
					return team;
				}),
				'Teams',
				false,
				() => {}
			)
		);
		hierarchy.push(getHierarchy('add_users', [], 'Users', true, () => {}, undefined, 'TAB'));
		setHeirarchy(hierarchy);
	}, [environmentsMap, org]);

	useEffect(() => {
		if (hierarchy.length === 0) {
			return;
		}

		if (!selectedHierarchy) {
			selectHierarchy(hierarchy[0]);
			return;
		} else if (selectedHierarchy.id === 'add_users') {
			selectHierarchy(hierarchy[2]);
		}
	}, [hierarchy]);

	return getListView(hierarchy, selectHierarchy, selectedHierarchy, 0);
};

const getListView = (
	hierarchies: Hierarchy[],
	selectHierarchy: (h: Hierarchy) => void,
	selectedHierarchy: Hierarchy | undefined,
	level: number
) => {
	const listItem = (hierarchy: Hierarchy) => {
		const expandable = hierarchy.children.length > 0;
		return (
			<li key={hierarchy.id} id={hierarchy.id}>
				<button
					onClick={() => {
						if (!hierarchy.selectable) return;
						selectHierarchy(hierarchy);
					}}
					className={`${hierarchy === selectedHierarchy ? 'selected' : ''} ${
						expandable ? 'expandable' : ''
					} ${hierarchy.expanded ? 'expanded' : ''}`}>
					{hierarchy.name}{' '}
					{expandable && (
						<span
							onClick={e => {
								hierarchy.expanded = e.currentTarget.parentElement?.classList.toggle('expanded');
							}}>
							<DropdownArrow />
						</span>
					)}
				</button>
				{expandable && getListView(hierarchy.children, selectHierarchy, selectedHierarchy, level + 1)}
			</li>
		);
	};

	return (
		<ul key={`secrets-container__list--${level}`} className={`secrets-container__list--${level}`}>
			{hierarchies.map(h => listItem(h))}
		</ul>
	);
};
