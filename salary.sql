SELECT 
  COALESCE(remarks, 'Total') AS Remarks, 
  FORMAT( sum(amount), 3) as totalSalary, 
  TIMESTAMPDIFF(MONTH, min(event_date), max(event_date)) as monthsServed
FROM `activities` as act
WHERE ( act.tag_id = 3 or act.sub_tag_id = 3 )
group by remarks WITH ROLLUP
order by max(event_date);