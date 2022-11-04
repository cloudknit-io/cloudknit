import React from 'react';
import { FC } from 'react';
import { ReactComponent as DownloadFile } from 'assets/images/icons/download.svg';


export type Props = {
	visualizationUrl: string;
    setVisualizationUrl: any;
};

export const Visualization: FC<Props> = ({visualizationUrl, setVisualizationUrl}: Props) => {
	let zoomLevel = 1;
	const visRef = React.useRef<HTMLDivElement>(null);

    setImmediate(() => {
        const svg = (visRef.current as HTMLDivElement).querySelector('.svg-dom > svg') as SVGElement;
        svg.style.transformOrigin = 'top left';
        svg.style.transition = 'transform 0.1s';
    });

    const downloadData = () => {
		if (!visualizationUrl) return;
		const blob = new Blob([visualizationUrl], { type: 'text' });
		const downloadLink = document.createElement('a');
		downloadLink.href = window.URL.createObjectURL(blob);
		downloadLink.download = 'visualization.svg';
		document.body.appendChild(downloadLink);
		downloadLink.click();
		document.body.removeChild(downloadLink);
	};

	return (
		<div ref={visRef} className="visualization">
			<div className="visualization_menu">
				<button
					onClick={() => {
						const svg = (visRef.current as HTMLDivElement).querySelector('.svg-dom > svg') as SVGElement;
						zoomLevel += 0.1;
						svg.style.transform = `scale(${zoomLevel})`;
					}}>
					+
				</button>
				<button
					onClick={() => {
						const svg = (visRef.current as HTMLDivElement).querySelector('.svg-dom > svg') as SVGElement;
						zoomLevel -= 0.1;
						svg.style.transform = `scale(${zoomLevel})`;
					}}>
					-
				</button>
				<button
					onClick={() => {
						downloadData();
					}}>
					<DownloadFile title="Download SVG" />
				</button>
				<button
					onClick={() => {
						setVisualizationUrl('');
					}}>
					x
				</button>
			</div>
			<div
				style={{ height: '100%', width: '100%', overflow: 'scroll' }}
				className={`svg-dom ${visualizationUrl ? '' : 'loading'}`}
				dangerouslySetInnerHTML={{ __html: visualizationUrl }}></div>
		</div>
	);
};
