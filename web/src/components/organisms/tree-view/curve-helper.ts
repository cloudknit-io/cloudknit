import * as d3 from 'd3';
export const getCurveType = (type: string): any => {
	switch (type) {
		case '0':
			return d3.curveBasis;
			break;
		case '1':
			return d3.curveBundle;
			break;
		case '2':
			return d3.curveCardinal;
			break;
		case '3':
			return d3.curveCatmullRom;
			break;
		case '4':
			return d3.curveLinear;
			break;
		case '5':
			return d3.curveMonotoneX;
			break;
		case '6':
			return d3.curveNatural;
			break;
		case '7':
			return d3.curveStep;
			break;
		default:
	}
};
