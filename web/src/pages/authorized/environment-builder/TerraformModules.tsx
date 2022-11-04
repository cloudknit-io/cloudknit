import React, { useEffect, useState } from 'react';
import ApiClient from 'utils/apiClient';

export const TerraformModules: React.FC = () => {
	const [modules, setModules] = useState<any[]>([]);
	const [allModules, setAllModules] = useState<any[]>([]);
	useEffect(() => {
		ApiClient.get('/terraform-external/modules/aws').then(({ data }) => {
			const all = (data as any).data.map((d: any) => ({
				fullName: d.attributes['full-name'],
				name: d.attributes['name'],
				source: d.attributes['source'],
				...d,
			}));
			setAllModules(all);
			setModules(all);
		});
	}, []);
	return (
		<>
			{/* <h3>Modules</h3> */}
			<div className="terraform-module">
				<input
					type="search"
					results={0}
					placeholder="Search modules"
					className="shadowy-input"
					onChange={ev => {
						const v = ev.currentTarget.value;
						if (v.trim() === '') {
							setModules(allModules);
						} else {
							setModules(allModules.filter(e => e.name.startsWith(v)));
						}
					}}
				/>
				<ul>
					{modules.map(m => (
						<li
							key={m.fullName}
							draggable="true"
							onDragStart={e => {
								e.dataTransfer.setData('text', JSON.stringify(m));
							}}>
							<label className="terraform-module__name terraform-module__name--bold">{m.name}</label>
							{/* <label className="terraform-module__fullname terraform-module__fullname--sub">
								{m.fullName}
							</label> */}
						</li>
					))}
				</ul>
			</div>
		</>
	);
};
