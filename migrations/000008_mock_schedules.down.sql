DELETE FROM schedules
WHERE group_id IN (1, 2, 3, 4)
  AND day_of_week IN (2, 3, 4, 5)
  AND time_of_day IN ('10:00', '15:00', '09:00', '14:00');