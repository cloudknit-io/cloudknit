import createNodeFigure, { updateNodeFigure } from './node-figure-helper';

export const getSVGNode = (
	attributes: { [key: string]: string },
	svgElementName: string,
	jsProps?: { [key: string]: any }
) => {
	const svgNode: any = document.createElementNS('http://www.w3.org/2000/svg', svgElementName);
	for (const key in attributes) {
		const value = attributes[key];
		svgNode.setAttribute(key, value);
	}

	for (const key in jsProps) {
		if (key) {
			const value = jsProps[key];
			svgNode[key] = value;
		}
	}
	return svgNode;
};

export const getShapeLabel = (item: any) => {
	const label = createNodeFigure({
		...item
	});

	return {
		label,
		labelType: 'svg',
		style: 'fill: transparent;',
		paddingX: 0,
		paddingY: 0,
	};
};

export const updateShapeLabel = (item: any) => {
	return updateNodeFigure({
		...item
	});
};
