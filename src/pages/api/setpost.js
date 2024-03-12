export async function post({request}) {
  const data = await request.formData(); // Here's the data sent by the form
  const text = data.get('text'); // Here's how you get the value of a field
  console.log(text);
  return new Response(JSON.stringify(request), {
    status: 200,
  });
}