import { useMemo, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export const useCalendarConfig = () => {
  const { t, i18n } = useTranslation();

  const calendarConfig = useMemo(() => {
    const currentLang = i18n.language;

    // Configuration de base pour le format 24h
    const baseConfig = {
      slotMinTime: '07:00:00',
      slotMaxTime: '22:00:00',
      slotDuration: '00:30:00',
      slotLabelFormat: {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false, // Format 24h
      },
      titleFormat: {
        year: 'numeric',
        month: 'long',    // Version longue pour le titre principal
        day: 'numeric',
      },
      buttonText: {
        today: t('calendar_today', 'Aujourd\'hui'),
        month: t('calendar_month', 'Mois'),
        week: t('calendar_week', 'Semaine'),
        day: t('calendar_day', 'Jour'),
        list: t('calendar_list', 'Liste'),
      },
      firstDay: 1, // Lundi comme premier jour de la semaine
      height: 'auto',
      expandRows: true,
      nowIndicator: true,
      businessHours: {
        daysOfWeek: [1, 2, 3, 4, 5], // Lundi à Vendredi
        startTime: '08:00',
        endTime: '18:00',
      },
    };

    // Configuration spécifique selon la langue
    const languageConfigs = {
      fr: {
        ...baseConfig,
        locale: 'fr',
        dayHeaderFormat: {
          weekday: 'short',
          month: 'short',
          day: 'numeric',
        },
        views: {
          dayGrid: {
            dayHeaderFormat: {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            },
          },
        },
      },
      en: {
        ...baseConfig,
        locale: 'en-gb', // Format 24h pour l'anglais
        dayHeaderFormat: {
          weekday: 'short',
          month: 'short',
          day: 'numeric',
        },
        views: {
          dayGrid: {
            dayHeaderFormat: {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            },
          },
        },
      },
      vi: {
        ...baseConfig,
        locale: 'vi',
        dayHeaderFormat: {
          weekday: 'short',
          month: 'short',
          day: 'numeric',
        },
        views: {
          dayGrid: {
            dayHeaderFormat: {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            },
          },
        },
      },
      cn: {
        ...baseConfig,
        locale: 'zh-cn',
        dayHeaderFormat: {
          weekday: 'short',
          month: 'short',
          day: 'numeric',
        },
        views: {
          dayGrid: {
            dayHeaderFormat: {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            },
          },
        },
      },
      ar: {
        ...baseConfig,
        locale: 'ar',
        dayHeaderFormat: {
          weekday: 'short',
          month: 'short',
          day: 'numeric',
        },
        views: {
          dayGrid: {
            dayHeaderFormat: {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            },
          },
        },
      },
    };

    return languageConfigs[currentLang] || languageConfigs.en;
  }, [t, i18n.language]);

  // Forcer la mise à jour quand la langue change
  useEffect(() => {
    // Cette fonction sera appelée à chaque changement de langue
    // Le useMemo ci-dessus se recalculera automatiquement
  }, [i18n.language]);

  return calendarConfig;
}; 