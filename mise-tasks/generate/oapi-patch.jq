"/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}/composites/clients/{client-uuid2}" as $new |
"/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}/composites/clients/{client-uuid}" as $old |
if .paths[$old] != null then
  .paths[$new] = .paths[$old] |
  del(.paths[$old]) |
  .paths[$new] |= with_entries(
    if (.value | type) == "object" and .value.parameters then
      .value.parameters = (.value.parameters + ($patch[0].paths[$new][.key].parameters // []))
    else . end
  )
else
  .
end
