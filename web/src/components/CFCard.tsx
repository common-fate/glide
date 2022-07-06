import { Box } from '@chakra-ui/layout';
import { BoxProps, forwardRef } from '@chakra-ui/react';

export const CFCard = forwardRef<BoxProps, 'div'>((props, ref) => (
	<Box
		p={5}
		borderRadius="18px"
		backgroundColor="white"
		boxShadow="0px 1px 4px rgba(91, 104, 113, 0.12)" // Figma: Shadow/01
		ref={ref}
		{...props}
	/>
));
