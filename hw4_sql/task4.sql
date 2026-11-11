SELECT DISTINCT name
FROM users
LEFT JOIN orders
ON orders.user_id = users.id
WHERE orders.item = 'наушники'
ORDER BY name;