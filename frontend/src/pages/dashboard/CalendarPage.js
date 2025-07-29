import FullCalendar from '@fullcalendar/react'; // => request placed at the top
import interactionPlugin from '@fullcalendar/interaction';
import listPlugin from '@fullcalendar/list';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import timelinePlugin from '@fullcalendar/timeline';
//
import { useState, useRef, useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import { Card, Container } from '@mui/material';
// hooks
import { useCalendar } from '../../hooks/useCalendar';
import { useCourses } from '../../hooks/useCourses';
import { useSubjects } from '../../hooks/useSubjects';
import { useTeachers } from '../../hooks/useTeachers';
import { useRooms } from '../../hooks/useRooms';
import { useCalendarConfig } from '../../hooks/useCalendarConfig';
// routes
import { PATH_DASHBOARD } from '../../routes/paths';
// utils
import { fTimestamp } from '../../utils/formatTime';
// hooks
import useResponsive from '../../hooks/useResponsive';
// components
import { useSnackbar } from '../../components/snackbar';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs';
import { useSettingsContext } from '../../components/settings';
import { useDateRangePicker } from '../../hooks/useDateRangePicker';
// sections
import {
  StyledCalendar,
  CalendarToolbar,
  CalendarFilterDrawer,
} from '../../sections/@dashboard/calendar';
import CourseEventDialog from '../../sections/@dashboard/calendar/CourseEventDialog';
import CourseFormDialog from '../../sections/courses/CourseFormDialog';

// ----------------------------------------------------------------------

export default function CalendarPage() {
  const { enqueueSnackbar } = useSnackbar();

  const { themeStretch } = useSettingsContext();

  const { events, updateEvent, deleteEvent } = useCalendar();
  const { createCourse, checkConflicts } = useCourses();
  const { fetchSubjects } = useSubjects();
  const { fetchTeachers } = useTeachers();
  const { fetchRooms } = useRooms();

  const isDesktop = useResponsive('up', 'sm');

  const calendarRef = useRef(null);

  const [openCourseDialog, setOpenCourseDialog] = useState(false);
  const [openFormDialog, setOpenFormDialog] = useState(false);

  const [selectedEventId, setSelectedEventId] = useState(null);
  const [initialFormData, setInitialFormData] = useState(null);

  const selectedEvent = selectedEventId ? events.find((event) => event.extendedProps?.course?.id === selectedEventId) : null;
  
  const picker = useDateRangePicker(null, null);

  const [date, setDate] = useState(new Date());

  const [openFilter, setOpenFilter] = useState(false);

  // Suppression du filtre de couleur

  const [view, setView] = useState(() => {
    // Initialiser la vue selon la taille de l'écran
    return isDesktop ? 'timeGridWeek' : 'listWeek';
  });

  // Configuration du calendrier selon la langue
  const calendarConfig = useCalendarConfig();

  useEffect(() => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();

      // Utiliser la vue semaine par défaut, sauf sur mobile où on utilise la liste
      const newView = isDesktop ? 'timeGridWeek' : 'listWeek';
      calendarApi.changeView(newView);
      setView(newView);
    }
  }, [isDesktop]);

  // Charger les données nécessaires pour le formulaire de cours
  useEffect(() => {
    fetchSubjects();
    fetchTeachers();
    fetchRooms();
  }, [fetchSubjects, fetchTeachers, fetchRooms]);

  const handleOpenCourseDialog = () => {
    setOpenCourseDialog(true);
  };

  const handleCloseCourseDialog = () => {
    setOpenCourseDialog(false);
    setSelectedEventId(null);
  };

  const handleCloseFormDialog = () => {
    setOpenFormDialog(false);
    setInitialFormData(null);
  };

  const handleSubmitCourse = async (data) => {
    try {
      const conflictsData = await checkConflicts(data);
      if (conflictsData && conflictsData.has_conflicts) {
        return conflictsData;
      }

      await createCourse(data);
      enqueueSnackbar('Cours créé avec succès!', { variant: 'success' });
      handleCloseFormDialog();
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleClickToday = () => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();
      calendarApi.today();
      setDate(calendarApi.getDate());
    }
  };

  const handleChangeView = (newView) => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();
      calendarApi.changeView(newView);
      setView(newView);
    }
  };

  const handleClickDatePrev = () => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();
      calendarApi.prev();
      setDate(calendarApi.getDate());
    }
  };

  const handleClickDateNext = () => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();
      calendarApi.next();
      setDate(calendarApi.getDate());
    }
  };

  const handleSelectRange = (arg) => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();
      calendarApi.unselect();
    }

    setInitialFormData({
      start_time: new Date(arg.start),
      end_time: new Date(arg.end),
    });
    setOpenFormDialog(true);
  };

  const handleSelectEvent = (arg) => {
    setSelectedEventId(arg.event.extendedProps?.course?.id);
    handleOpenCourseDialog();
  };

  const handleResizeEvent = async ({ event }) => {
    try {
      const updatedEvent = {
        title: event.title,
        start: event.start,
        end: event.end,
        allDay: event.allDay,
      };

      await updateEvent(event.id, updatedEvent);
      enqueueSnackbar('Cours mis à jour avec succès!', { variant: 'success' });
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleDropEvent = async ({ event }) => {
    try {
      const updatedEvent = {
        title: event.title,
        start: event.start,
        end: event.end,
        allDay: event.allDay,
      };

      await updateEvent(event.id, updatedEvent);
      enqueueSnackbar('Cours mis à jour avec succès!', { variant: 'success' });
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleCreateUpdateEvent = async (newEvent) => {
    try {
      if (selectedEventId) {
        await updateEvent(selectedEventId, newEvent);
        enqueueSnackbar('Cours mis à jour avec succès!', { variant: 'success' });
      } else {
        window.location.href = '/dashboard/courses';
        enqueueSnackbar('Redirection vers la création de cours...', { variant: 'info' });
      }
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleDeleteEvent = async () => {
    try {
      const event = events.find((e) => e.id === selectedEventId);
      if (event && event.extendedProps.course) {
        handleCloseCourseDialog();
        await deleteEvent(event.extendedProps.course.id);
        enqueueSnackbar('Cours supprimé avec succès!', { variant: 'success' });
      }
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleResetFilter = () => {
    const { setStartDate, setEndDate } = picker;

    if (setStartDate && setEndDate) {
      setStartDate(null);
      setEndDate(null);
    }
  };

  // Transformer les événements pour inclure textColor
  const transformedEvents = events.map((event) => ({
    ...event,
    textColor: event.color,
  }));

  const dataFiltered = applyFilter({
    inputData: transformedEvents,
    filterStartDate: picker.startDate,
    filterEndDate: picker.endDate,
    isError: !!picker.isError,
  });

  return (
    <>
      <Helmet>
        <title> Calendrier | EduQR</title>
      </Helmet>

      <Container maxWidth={themeStretch ? false : 'xl'}>
        <CustomBreadcrumbs
          heading="Calendrier"
          links={[
            {
              name: 'Dashboard',
              href: PATH_DASHBOARD.root,
            },
            {
              name: 'Calendrier',
            },
          ]}
          moreLink={['https://fullcalendar.io/docs/react']}
        />

        <Card>
          <StyledCalendar>
            <CalendarToolbar
              date={date}
              view={view}
              onNextDate={handleClickDateNext}
              onPrevDate={handleClickDatePrev}
              onToday={handleClickToday}
              onChangeView={handleChangeView}
              onOpenFilter={() => setOpenFilter(true)}
              onAddNew={() => {
                // Ouvrir la modal de création de cours
                setInitialFormData(null);
                setOpenFormDialog(true);
              }}
            />

            <FullCalendar
              key={`calendar-${calendarConfig.locale}`}
              {...calendarConfig}
              weekends
              editable
              droppable
              selectable
              rerenderDelay={10}
              allDayMaintainDuration
              eventResizableFromStart
              ref={calendarRef}
              initialDate={date}
              initialView={view}
              dayMaxEventRows={3}
              eventDisplay="block"
              events={dataFiltered}
              headerToolbar={false}
              select={handleSelectRange}
              eventDrop={handleDropEvent}
              eventClick={handleSelectEvent}
              eventResize={handleResizeEvent}
              height={isDesktop ? 720 : 'auto'}
              plugins={[
                listPlugin,
                dayGridPlugin,
                timelinePlugin,
                timeGridPlugin,
                interactionPlugin,
              ]}
              views={{
                dayGrid: {
                  dayHeaderFormat: false, // Supprime les en-têtes de jours dans la vue mensuelle
                },
              }}
            />
          </StyledCalendar>
        </Card>
      </Container>

      <CalendarFilterDrawer
        events={events}
        picker={picker}
        openFilter={openFilter}
        onResetFilter={handleResetFilter}
        onCloseFilter={() => setOpenFilter(false)}
        onSelectEvent={(eventId) => {
          if (eventId) {
            setSelectedEventId(eventId);
            handleOpenCourseDialog();
          }
        }}
      />

      <CourseEventDialog
        open={openCourseDialog}
        onClose={handleCloseCourseDialog}
        event={selectedEvent}
        onUpdate={handleCreateUpdateEvent}
        onDelete={handleDeleteEvent}
      />

      <CourseFormDialog
        open={openFormDialog}
        onClose={handleCloseFormDialog}
        onSubmit={handleSubmitCourse}
        initialData={initialFormData}
      />
    </>
  );
}

function applyFilter({ inputData, filterStartDate, filterEndDate, isError }) {
  if (filterStartDate && filterEndDate && !isError) {
    return inputData.filter(
      (event) =>
        fTimestamp(event.start) >= fTimestamp(filterStartDate) &&
        fTimestamp(event.end) <= fTimestamp(filterEndDate)
    );
  }

  return inputData;
} 