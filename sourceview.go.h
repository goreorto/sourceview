#include <gtk/gtk.h>

static inline gchar** make_strings(int count) {
	return (gchar**)malloc(sizeof(gchar*) * count);
}

static inline void destroy_strings(gchar** strings) {
	free(strings);
}

static inline gchar* get_string(gchar** strings, int n) {
	return strings[n];
}

static inline void set_string(gchar** strings, int n, gchar* str) {
	strings[n] = str;
}

static inline gchar** next_gcharptr(gchar** s) { return (s+1); }

static const gchar *
object_get_class_name(GObject *object)
{
	return G_OBJECT_CLASS_NAME(G_OBJECT_GET_CLASS(object));
}

static GObject *
toGObject(void *p)
{
	return (G_OBJECT(p));
}

static GtkTextView *
toGtkTextView(void *p)
{
	return (GTK_TEXT_VIEW(p));
}

static GtkTextBuffer *
toGtkTextBuffer(void *p)
{
	return (GTK_TEXT_BUFFER(p));
}

static GtkSourceView *
toGtkSourceView(void *p)
{
	return (GTK_SOURCE_VIEW(p));
}

static GtkSourceBuffer *
toGtkSourceBuffer(void *p)
{
	return (GTK_SOURCE_BUFFER(p));
}

static GtkSourceCompletion *
toGtkSourceCompletion(void *p)
{
	return (GTK_SOURCE_COMPLETION(p));
}

static GtkSourceCompletionContext *
toGtkSourceCompletionContext(void *p)
{
	return (GTK_SOURCE_COMPLETION_CONTEXT(p));
}

static GtkSourceCompletionInfo *
toGtkSourceCompletionInfo(void *p)
{
	return (GTK_SOURCE_COMPLETION_INFO(p));
}

static GtkSourceCompletionProposal *
toGtkSourceCompletionProposal(void *p)
{
	return (GTK_SOURCE_COMPLETION_PROPOSAL(p));
}

static GtkSourceCompletionProvider *
toGtkSourceCompletionProvider(void *p)
{
	return (GTK_SOURCE_COMPLETION_PROVIDER(p));
}

static GtkSourceGutter *
toGtkSourceGutter(void *p)
{
	return (GTK_SOURCE_GUTTER(p));
}

static GtkSourceLanguageManager *
toGtkSourceLanguageManager(void *p)
{
	return (GTK_SOURCE_LANGUAGE_MANAGER(p));
}

static GtkSourceLanguage *
toGtkSourceLanguage(void *p)
{
	return (GTK_SOURCE_LANGUAGE(p));
}

static GtkSourceStyle *
toGtkSourceStyle(void *p)
{
	return (GTK_SOURCE_STYLE(p));
}

static GtkSourceStyleScheme *
toGtkSourceStyleScheme(void *p)
{
	return (GTK_SOURCE_STYLE_SCHEME(p));
}

static GtkSourceStyleSchemeManager *
toGtkSourceStyleSchemeManager(void *p)
{
	return (GTK_SOURCE_STYLE_SCHEME_MANAGER(p));
}

static GtkSourceStyleSchemeChooser *
toGtkSourceStyleSchemeChooser(void *p)
{
	return (GTK_SOURCE_STYLE_SCHEME_CHOOSER(p));
}

static GtkSourceStyleSchemeChooserButton *
toSourceStyleSchemeChooserButton(void *p)
{
	return (GTK_SOURCE_STYLE_SCHEME_CHOOSER_BUTTON(p));
}

static GtkSourceStyleSchemeChooserWidget *
toSourceStyleSchemeChooserWidget(void *p)
{
	return (GTK_SOURCE_STYLE_SCHEME_CHOOSER_WIDGET(p));
}

static GtkTextTag *
toGtkTextTag(void *p)
{
	return (GTK_TEXT_TAG(p));
}
