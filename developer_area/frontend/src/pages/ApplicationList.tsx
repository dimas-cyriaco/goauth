import { For, createResource } from 'solid-js'

import { application } from '../lib/client'
import { createAPIClient } from '../lib/clientUtils'

export const ApplicationList = () => {
  const [applications] = createResource(fetchApplications)

  return (
    <div>
      <h1>Applications</h1>

      <p>List of all applications</p>

      <For each={applications()} fallback={<div>Loading...</div>}>
        {(item) => <div data-testid="application-name">{item.name}</div>}
      </For>
    </div>
  )
}

const fetchApplications = async () => {
  try {
    const client = createAPIClient()

    const params: application.ApplicationListParams = {
      Page: 1,
      PerPage: 10,
    }

    const response = await client.application.List(params)
    return response.applications
  } catch {
    return []
  }
}
