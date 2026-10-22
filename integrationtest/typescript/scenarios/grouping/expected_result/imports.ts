import type { X } from 'hellyeah';

import { redirect } from '@sveltejs/kit';
import * as R from 'ramda';
import Rand from 'rand-seed';

import { add } from '../helpers/maths';
import { randomCoordinates } from '../helpers/random';
import type { Colour, Line, Square } from '../helpers/rendering';
import { blur } from '../helpers/rendering';

import type { Y } from './but';

import { AvailabilityStatus, Availability, AvailabilityDatetime } from './lib/types/availability';

const x = 1;


function foo() {
	return 1;
}

