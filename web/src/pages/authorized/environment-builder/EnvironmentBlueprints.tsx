import { environmentBlueprints } from 'helpers/environment-builder.helper';
import React, { useEffect, useState } from 'react';

export const EnvironmentTemplates: React.FC = () => {
	const [blueprints, setBlueprints] = useState<any[]>([]);

	useEffect(() => {
		setBlueprints(environmentBlueprints)
	}, []);
	return (
		<>
			<div className="templates">
				<input
					type="search"
					results={0}
					placeholder="Search Blueprints"
					className="shadowy-input"
					onChange={ev => {
						const v = ev.currentTarget.value;
						if (v.trim() === '') {
							setBlueprints(environmentBlueprints);
						} else {
							setBlueprints(blueprints.filter(e => e.id.toLowerCase().includes(v.toLowerCase())));
						}
					}}
				/>
				<ul>
					{blueprints.map(b => (
						<li
							key={b.id}
							draggable="true"
							onDragStart={e => {
								e.dataTransfer.setData('text', JSON.stringify({type: 'template' , id: b.id }));
							}}>
							<label className="terraform-module__name terraform-module__name--bold">{b.id}</label>
						</li>
					))}
				</ul>
			</div>
		</>
	);
};
