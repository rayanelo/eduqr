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
import { Card, Button, Container } from '@mui/material';
// hooks
import { useCalendar } from '../../hooks/useCalendar';
import { useCourses } from '../../hooks/useCourses';
import { useSubjects } from '../../hooks/useSubjects';
import { useTeachers } from '../../hooks/useTeachers';
import { useRooms } from '../../hooks/useRooms';
// routes
import { PATH_DASHBOARD } from '../../routes/paths';
// utils
import { fTimestamp } from '../../utils/formatTime';
// hooks
import useResponsive from '../../hooks/useResponsive';
// components
import Iconify from '../../components/iconify';
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

// Suppression des options de couleur - tous les cours auront la même couleur

// ----------------------------------------------------------------------

export default function CalendarPage() {
  const { enqueueSnackbar } = useSnackbar();

  const { themeStretch } = useSettingsContext();

  const { events, updateEvent, deleteEvent } = useCalendar();
  const { createCourse, updateCourse, deleteCourse, checkConflicts } = useCourses();
  const { subjects, fetchSubjects } = useSubjects();
  const { teachers, fetchTeachers } = useTeachers();
  const { rooms, fetchRooms } = useRooms();

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

  const [view, setView] = useState(isDesktop ? 'dayGridMonth' : 'listWeek');

  useEffect(() => {
    const calendarEl = calendarRef.current;
    if (calendarEl) {
      const calendarApi = calendarEl.getApi();

      const newView = isDesktop ? 'dayGridMonth' : 'listWeek';
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
      // Vérifier les conflits avant de créer
      const conflictsData = await checkConflicts(data);
      if (conflictsData && conflictsData.has_conflicts) {
        return conflictsData; // Retourner les conflits pour affichage
      }

      await createCourse(data);
      enqueueSnackbar('Cours créé avec succès!', { variant: 'success' });
      handleCloseFormDialog();
      
      // Recharger les événements du calendrier
      window.location.reload(); // Solution temporaire pour recharger
    } catch (error) {
      console.error('Erreur lors de la création du cours:', error);
      enqueueSnackbar('Erreur lors de la création du cours', { variant: 'error' });
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
    // Ouvrir la modal de création de cours avec les dates pré-remplies
    setInitialFormData({
      start_time: new Date(arg.start),
      end_time: new Date(arg.end),
    });
    setOpenFormDialog(true);
  };

  const handleSelectEvent = (arg) => {
    // Utiliser l'ID du cours depuis extendedProps au lieu de l'ID de FullCalendar
    const courseId = arg.event.extendedProps?.course?.id;
    setSelectedEventId(courseId);
    // Tous les événements sont des cours maintenant
    handleOpenCourseDialog();
  };

  const handleResizeEvent = async ({ event }) => {
    try {
      await updateEvent(event.id, {
        allDay: event.allDay,
        start: event.start,
        end: event.end,
      });
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Error updating event!', { variant: 'error' });
    }
  };

  const handleDropEvent = async ({ event }) => {
    try {
      await updateEvent(event.id, {
        allDay: event.allDay,
        start: event.start,
        end: event.end,
      });
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Error updating event!', { variant: 'error' });
    }
  };

  const handleCreateUpdateEvent = async (newEvent) => {
    try {
      if (selectedEventId) {
        await updateEvent(selectedEventId, newEvent);
        enqueueSnackbar('Cours mis à jour avec succès!');
      } else {
        // Rediriger vers la création de cours
        window.location.href = '/dashboard/courses';
        enqueueSnackbar('Redirection vers la création de cours...');
      }
    } catch (error) {
      console.error(error);
      enqueueSnackbar('Une erreur est survenue!', { variant: 'error' });
    }
  };

  const handleDeleteEvent = async () => {
    try {
      if (selectedEventId) {
        // Trouver l'événement correspondant pour obtenir l'ID du cours
        const event = events.find(e => e.id === selectedEventId);
        if (event && event.extendedProps.course) {
          handleCloseCourseDialog();
          await deleteEvent(event.extendedProps.course.id);
          enqueueSnackbar('Cours supprimé avec succès!');
        }
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
          action={
            <Button
              variant="contained"
              startIcon={<Iconify icon="eva:plus-fill" />}
              onClick={() => {
                setInitialFormData(null);
                setOpenFormDialog(true);
              }}
            >
              Nouveau Cours
            </Button>
          }
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
  const stabilizedThis = inputData.map((el, index) => [el, index]);

  inputData = stabilizedThis.map((el) => el[0]);

  if (filterStartDate && filterEndDate && !isError) {
    inputData = inputData.filter(
      (event) =>
        fTimestamp(event.start) >= fTimestamp(filterStartDate) &&
        fTimestamp(event.end) <= fTimestamp(filterEndDate)
    );
  }

  return inputData;
} 