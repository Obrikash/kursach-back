DELETE FROM user_subscriptions
WHERE (user_id = 1 AND subscription_id = 1)
   OR (user_id = 2 AND subscription_id = 2)
   OR (user_id = 3 AND subscription_id = 3)
   OR (user_id = 4 AND subscription_id = 4);