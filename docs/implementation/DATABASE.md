# Database Implementation Notes

Source of truth for schema:
- `middleware/database/schema.go`

If you change any table/column/index:
- Update `middleware/database/schema.go`
- Consider migration scripts (if needed)
- Update this file with brief rationale

Recent schema additions:
- `coops.temp_min` / `coops.temp_max` for threshold control
- `coops.water_level_half_threshold` for water alert calibration
- `alerts.is_acknowledged` for alert state
- `devices.hardware_id` is no longer unique (multiple devices per gateway)
