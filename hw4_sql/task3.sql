SELECT name, COALESCE(SUM(price), 0) as total_price
FROM users
LEFT JOIN orders
ON users.id = order.user_id
GROUP BY users.id
ORDER BY COALESCE(SUM(price), 0) DESC;