import _mock from '../_mock';
import { addDays, subDays, addHours, addMinutes } from 'date-fns';

// ----------------------------------------------------------------------

const now = new Date();

export const _calendarEvents = [...Array(24)].map((_, index) => {
  const setIndex = index + 1;

  const start = addDays(subDays(now, 12), Math.floor(index / 2));

  const end = addHours(addMinutes(start, Math.floor(index / 5) * 15), Math.floor(index / 2));

  return {
    id: _mock.id(setIndex),
    title: (() => {
      const titles = [
        'Meeting with Client',
        'Team Standup',
        'Project Review',
        'Design Workshop',
        'Code Review',
        'Product Demo',
        'Strategy Meeting',
        'Client Presentation',
        'Sprint Planning',
        'Retrospective',
        'User Testing',
        'Deployment',
        'Training Session',
        'Interview',
        'Conference Call',
        'Workshop',
        'Brainstorming',
        'Status Update',
        'Planning Session',
        'Review Meeting',
        'Kickoff Meeting',
        'Follow-up Call',
        'Technical Discussion',
        'Stakeholder Meeting',
      ];
      return titles[index % titles.length];
    })(),
    description: _mock.text.description(setIndex),
    color: (() => {
      const colors = [
        '#00AB55', // green
        '#1890FF', // blue
        '#54D62C', // light green
        '#FFC107', // yellow
        '#FF4842', // red
        '#04297A', // dark blue
        '#7A0C2E', // dark red
      ];
      return colors[index % colors.length];
    })(),
    allDay: index % 3 === 0,
    start,
    end,
  };
});

export const _calendarCategories = [
  { id: 1, name: 'Work', color: '#00AB55' },
  { id: 2, name: 'Personal', color: '#1890FF' },
  { id: 3, name: 'Meeting', color: '#54D62C' },
  { id: 4, name: 'Important', color: '#FFC107' },
  { id: 5, name: 'Urgent', color: '#FF4842' },
  { id: 6, name: 'Project', color: '#04297A' },
  { id: 7, name: 'Review', color: '#7A0C2E' },
]; 