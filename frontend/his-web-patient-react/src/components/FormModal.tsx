import { Modal, Form, Button, Message } from 'semantic-ui-react'
import { useState, useEffect, type ReactNode } from 'react'

interface FormField {
  key: string
  label: string
  required?: boolean
  render: (value: string, onChange: (val: string) => void) => ReactNode
}

interface FormModalProps {
  open: boolean
  title: string
  fields: FormField[]
  initialValues: Record<string, string>
  onSubmit: (values: Record<string, string>) => Promise<void>
  onClose: () => void
}

export default function FormModal({ open, title, fields, initialValues, onSubmit, onClose }: FormModalProps) {
  const [values, setValues] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (open) {
      setValues({ ...initialValues })
      setError('')
    }
  }, [open, initialValues])

  const handleSubmit = async () => {
    // Check required fields
    for (const f of fields) {
      if (f.required && !values[f.key]) {
        setError(`请填写${f.label}`)
        return
      }
    }
    setLoading(true)
    setError('')
    try {
      await onSubmit(values)
      onClose()
    } catch (err: any) {
      setError(err.message || '操作失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} size="small" closeIcon>
      <Modal.Header>{title}</Modal.Header>
      <Modal.Content>
        {error && (
          <Message negative onDismiss={() => setError('')}>
            {error}
          </Message>
        )}
        <Form loading={loading}>
          {fields.map((field) => (
            <Form.Field key={field.key} required={field.required}>
              <label>{field.label}</label>
              {field.render(values[field.key] || '', (val) =>
                setValues((prev) => ({ ...prev, [field.key]: val })),
              )}
            </Form.Field>
          ))}
        </Form>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={onClose}>取消</Button>
        <Button primary loading={loading} onClick={handleSubmit}>
          确认
        </Button>
      </Modal.Actions>
    </Modal>
  )
}
