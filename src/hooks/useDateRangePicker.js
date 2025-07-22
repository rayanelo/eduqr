import { useState, useCallback } from 'react';
import { isBefore } from 'date-fns';

// ----------------------------------------------------------------------

export const useDateRangePicker = (startDate, endDate) => {
  const [open, setOpen] = useState(false);
  const [startDateValue, setStartDateValue] = useState(startDate);
  const [endDateValue, setEndDateValue] = useState(endDate);

  const isError = startDateValue && endDateValue && isBefore(new Date(endDateValue), new Date(startDateValue));

  const onChangeStartDate = useCallback((newValue) => {
    setStartDateValue(newValue);
  }, []);

  const onChangeEndDate = useCallback((newValue) => {
    setEndDateValue(newValue);
  }, []);

  const onOpen = useCallback(() => {
    setOpen(true);
  }, []);

  const onClose = useCallback(() => {
    setOpen(false);
  }, []);

  const onReset = useCallback(() => {
    setStartDateValue(null);
    setEndDateValue(null);
  }, []);

  return {
    startDate: startDateValue,
    endDate: endDateValue,
    onChangeStartDate,
    onChangeEndDate,
    isError,
    open,
    onOpen,
    onClose,
    onReset,
  };
}; 