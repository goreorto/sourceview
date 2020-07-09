package sourceview

// #cgo pkg-config: gtksourceview-3.0
// #include <gtksourceview/gtksourcebuffer.h>
// #include <gtksourceview/gtksourcecompletion.h>
// #include <gtksourceview/gtksourcecompletioncontext.h>
// #include <gtksourceview/gtksourcecompletioninfo.h>
// #include <gtksourceview/gtksourcecompletionproposal.h>
// #include <gtksourceview/gtksourcecompletionprovider.h>
// #include <gtksourceview/gtksourcegutter.h>
// #include <gtksourceview/gtksourcelanguage.h>
// #include <gtksourceview/gtksourcelanguagemanager.h>
// #include <gtksourceview/gtksourcestyle.h>
// #include <gtksourceview/gtksourcestylescheme.h>
// #include <gtksourceview/gtksourcestyleschemechooser.h>
// #include <gtksourceview/gtksourcestyleschemechooserbutton.h>
// #include <gtksourceview/gtksourcestyleschemechooserwidget.h>
// #include <gtksourceview/gtksourcestyleschememanager.h>
// #include <gtksourceview/gtksourceview.h>
// #include "sourceview.go.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var errNilPtr = errors.New("cgo returned unexpected nil pointer")

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gtk_source_buffer_get_type()), marshalSourceBuffer},
		{glib.Type(C.gtk_source_completion_get_type()), marshalSourceCompletion},
		{glib.Type(C.G_TYPE_ENUM), marshalSourceCompletionActivation},
		{glib.Type(C.gtk_source_completion_context_get_type()), marshalSourceCompletionContext},
		{glib.Type(C.gtk_source_completion_info_get_type()), marshalSourceCompletionInfo},
		{glib.Type(C.gtk_source_completion_proposal_get_type()), marshalSourceCompletionProposal},
		{glib.Type(C.gtk_source_completion_provider_get_type()), marshalSourceCompletionProvider},
		{glib.Type(C.gtk_source_gutter_get_type()), marshalSourceGutter},
		{glib.Type(C.gtk_source_language_get_type()), marshalSourceLanguage},
		{glib.Type(C.gtk_source_language_manager_get_type()), marshalSourceLanguageManager},
		{glib.Type(C.gtk_source_style_get_type()), marshalSourceStyle},
		{glib.Type(C.gtk_source_style_scheme_get_type()), marshalSourceStyleScheme},
		{glib.Type(C.gtk_source_style_scheme_chooser_get_type()), marshalSourceStyleSchemeChooser},
		{glib.Type(C.gtk_source_style_scheme_chooser_button_get_type()), marshalSourceStyleSchemeChooserButton},
		{glib.Type(C.gtk_source_style_scheme_chooser_widget_get_type()), marshalSourceStyleSchemeChooserWidget},
		{glib.Type(C.gtk_source_style_scheme_manager_get_type()), marshalSourceStyleSchemeManager},
		{glib.Type(C.gtk_source_view_get_type()), marshalSourceView},
	}
	glib.RegisterGValueMarshalers(tm)

	gtk.WrapMap["GtkSourceView"] = wrapSourceView
	gtk.WrapMap["GtkSourceBuffer"] = wrapSourceBuffer
	gtk.WrapMap["GtkSourceCompletion"] = wrapSourceCompletion
	gtk.WrapMap["GtkSourceCompletionContext"] = wrapSourceCompletionContext
	gtk.WrapMap["GtkSourceCompletionInfo"] = wrapSourceCompletionInfo
	gtk.WrapMap["GtkSourceCompletionProposal"] = wrapSourceCompletionProposal
	gtk.WrapMap["GtkSourceCompletionProvider"] = wrapSourceCompletionProvider
	gtk.WrapMap["GtkSourceGutter"] = wrapSourceGutter
	gtk.WrapMap["GtkSourceLanguage"] = wrapSourceLanguage
	gtk.WrapMap["GtkSourceLanguageManager"] = wrapSourceLanguageManager
	gtk.WrapMap["GtkSourceStyle"] = wrapSourceStyle
	gtk.WrapMap["GtkSourceStyleScheme"] = wrapSourceStyleScheme
	gtk.WrapMap["GtkSourceStyleSchemeChooser"] = wrapSourceStyleSchemeChooser
	gtk.WrapMap["GtkSourceStyleSchemeChooserButton"] = wrapSourceStyleSchemeChooserButton
	gtk.WrapMap["GtkSourceStyleSchemeChooserWidget"] = wrapSourceStyleSchemeChooserWidget
	gtk.WrapMap["GtkSourceStyleSchemeManager"] = wrapSourceStyleSchemeManager
}

func gobool(b C.gboolean) bool {
	return b != C.FALSE
}

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func goString(cstr *C.gchar) string {
	return C.GoString((*C.char)(cstr))
}

// TODO: cast* functions taken gotk3
func castWidget(c *C.GtkWidget) (gtk.IWidget, error) {
	ptr := unsafe.Pointer(c)
	var (
		className = goString(C.object_get_class_name(C.toGObject(ptr)))
		obj       = glib.Take(ptr)
	)

	intf, err := castInternal(className, obj)
	if err != nil {
		return nil, err
	}

	ret, ok := intf.(gtk.IWidget)
	if !ok {
		return nil, fmt.Errorf("expected value of type IWidget, got %T", intf)
	}

	return ret, nil
}

func castInternal(className string, obj *glib.Object) (interface{}, error) {
	fn, ok := gtk.WrapMap[className]
	if !ok {
		return nil, errors.New("unrecognized class name '" + className + "'")
	}

	// Check that the wrapper function is actually a function
	rf := reflect.ValueOf(fn)
	if rf.Type().Kind() != reflect.Func {
		return nil, errors.New("wraper is not a function")
	}

	// Call the wraper function with the *glib.Object as first parameter
	// e.g. "wrapWindow(obj)"
	v := reflect.ValueOf(obj)
	rv := rf.Call([]reflect.Value{v})

	// At most/max 1 return value
	if len(rv) != 1 {
		return nil, errors.New("wrapper did not return")
	}

	// Needs to be a pointer of some sort
	if k := rv[0].Kind(); k != reflect.Ptr {
		return nil, fmt.Errorf("wrong return type %s", k)
	}

	// Only get an interface value, type check will be done in more specific functions
	return rv[0].Interface(), nil
}

// SourceCompletionActivation is a representation of GTK's GtkSourceCompletionActivation.
type SourceCompletionActivation int

const (
	SOURCE_COMPLETION_ACTIVATION_NONE           SourceCompletionActivation = C.GTK_SOURCE_COMPLETION_ACTIVATION_NONE
	SOURCE_COMPLETION_ACTIVATION_INTERACTIVE    SourceCompletionActivation = C.GTK_SOURCE_COMPLETION_ACTIVATION_INTERACTIVE
	SOURCE_COMPLETION_ACTIVATION_USER_REQUESTED SourceCompletionActivation = C.GTK_SOURCE_COMPLETION_ACTIVATION_USER_REQUESTED
)

func marshalSourceCompletionActivation(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return SourceCompletionActivation(c), nil
}

/*
 * GtkSourceCompletion
 */

// SourceCompletion is a representation of GTK's GtkSourceCompletion.
type SourceCompletion struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceCompletion.
func (v *SourceCompletion) native() *C.GtkSourceCompletion {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceCompletion(p)
}

func marshalSourceCompletion(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceCompletion(obj), nil
}

func wrapSourceCompletion(obj *glib.Object) *SourceCompletion {
	return &SourceCompletion{obj}
}

// AddProvider is a wrapper around gtk_source_completion_add_provider().
func (v *SourceCompletion) AddProvider(provider ISourceCompletionProvider) error {
	var err *C.GError = nil
	cbool := C.gtk_source_completion_add_provider(v.native(), provider.toSourceCompletionProvider(), &err)
	if !gobool(cbool) {
		defer C.g_error_free(err)
		return errors.New(goString(err.message))
	}

	return nil
}

// RemoveProvider is a wrapper around gtk_source_completion_remove_provider().
func (v *SourceCompletion) RemoveProvider(provider ISourceCompletionProvider) error {
	var err *C.GError = nil
	cbool := C.gtk_source_completion_add_provider(v.native(), provider.toSourceCompletionProvider(), &err)
	if !gobool(cbool) {
		defer C.g_error_free(err)
		return errors.New(goString(err.message))
	}

	return nil
}

// GetProviders is a wrapper around gtk_source_completion_get_provider().
func (v *SourceCompletion) GetProviders() *glib.List {
	clist := C.gtk_source_completion_get_providers(v.native())
	if clist == nil {
		return nil
	}

	glist := glib.WrapList(uintptr(unsafe.Pointer(clist)))
	glist.DataWrapper(func(ptr unsafe.Pointer) interface{} {
		return &SourceCompletionProvider{glib.Take(ptr)}
	})

	return glist
}

// Show is a wrapper around gtk_source_completion_show().
func (v *SourceCompletion) Show(providers *glib.List, context *SourceCompletionContext) bool {
	nativeList := (*C.struct__GList)(unsafe.Pointer(providers.Native()))
	return gobool(C.gtk_source_completion_show(v.native(), nativeList, context.native()))
}

// Hide is a wrapper around gtk_source_completion_hide().
func (v *SourceCompletion) Hide() {
	C.gtk_source_completion_hide(v.native())
}

// GetInfoWindow is a wrapper around gtk_source_completion_get_info_window().
func (v *SourceCompletion) GetInfoWindow() (*SourceCompletionInfo, error) {
	c := C.gtk_source_completion_get_info_window(v.native())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceCompletionInfo(glib.Take(unsafe.Pointer(c))), nil
}

// GetView is a wrapper around gtk_source_completion_get_view().
func (v *SourceCompletion) GetView() (*SourceView, error) {
	c := C.gtk_source_completion_get_view(v.native())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceView(glib.Take(unsafe.Pointer(c))), nil
}

// CreateContext is a wrapper around gtk_source_completion_create_context().
func (v *SourceCompletion) CreateContext(position *gtk.TextIter) (*SourceCompletionContext, error) {
	// TODO: no idea if (*C.GtkTextIter)(unsafe.Pointer(position)) works
	c := C.gtk_source_completion_create_context(v.native(), (*C.GtkTextIter)(unsafe.Pointer(position)))
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceCompletionContext(glib.Take(unsafe.Pointer(c))), nil
}

// BlockInteractive is a wrapper around gtk_source_completion_block_interactive().
func (v *SourceCompletion) BlockInteractive() {
	C.gtk_source_completion_block_interactive(v.native())
}

// UnlockInteractive is a wrapper around gtk_source_completion_unblock_interactive().
func (v *SourceCompletion) UnlockInteractive() {
	C.gtk_source_completion_unblock_interactive(v.native())
}

/*
 * GtkSourceCompletionContext
 */

// SourceCompletionContext is a representation of GTK's GtkSourceCompletionContext.
type SourceCompletionContext struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceCompletionContext.
func (v *SourceCompletionContext) native() *C.GtkSourceCompletionContext {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceCompletionContext(p)
}

func marshalSourceCompletionContext(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceCompletionContext(obj), nil
}

func wrapSourceCompletionContext(obj *glib.Object) *SourceCompletionContext {
	return &SourceCompletionContext{obj}
}

// AddProposals is a wrapper around gtk_source_completion_context_add_proposals().
func (v *SourceCompletionContext) AddProposals(provider *SourceCompletionProvider, proposals *glib.List, finished bool) {
	nativeList := (*C.struct__GList)(unsafe.Pointer(proposals.Native()))
	C.gtk_source_completion_context_add_proposals(v.native(), provider.native(), nativeList, gbool(finished))
}

// GetIter is a wrapper around gtk_source_completion_context_get_iter().
func (v *SourceCompletionContext) GetIter() *gtk.TextIter {
	var iter C.GtkTextIter
	C.gtk_source_completion_context_get_iter(v.native(), &iter)
	// TODO: no idea if this casting is working
	return (*gtk.TextIter)(unsafe.Pointer(&iter))
}

// GetActivation is a wrapper around gtk_source_completion_get_activation().
func (v *SourceCompletionContext) GetActivation() SourceCompletionActivation {
	return SourceCompletionActivation(C.gtk_source_completion_context_get_activation(v.native()))
}

/*
 * GtkSourceCompletionInfo
 */

// SourceCompletionInfo is a representation of GTK's GtkSourceCompletionInfo.
type SourceCompletionInfo struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceCompletionInfo.
func (v *SourceCompletionInfo) native() *C.GtkSourceCompletionInfo {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceCompletionInfo(p)
}

func marshalSourceCompletionInfo(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceCompletionInfo(obj), nil
}

func wrapSourceCompletionInfo(obj *glib.Object) *SourceCompletionInfo {
	return &SourceCompletionInfo{obj}
}

// SourceCompletionInfoNew is a wrapper around gtk_source_completion_info_new().
func SourceCompletionInfoNew() (*SourceCompletionInfo, error) {
	c := C.gtk_source_completion_info_new()
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceCompletionInfo(glib.Take(unsafe.Pointer(c))), nil
}

// MoveToIter is a wrapper around gtk_source_completion_info_move_to_iter().
func (v *SourceCompletionInfo) MoveToIter(view *SourceView, iter *gtk.TextIter) {
	C.gtk_source_completion_info_move_to_iter(
		v.native(),
		view.asTextView(),
		(*C.GtkTextIter)(unsafe.Pointer(iter)),
	)
}

/*
 * GtkSourceCompletionProvider
 */

// ISourceCompletionProvider is a representation of the GtkSourceCompletionProvider GInterface,
// used to avoid duplication when embedding the type in a wrapper of another GObject-based type.
// The non-Interface version should only be used Actionable is used if the concrete type is not known.
type ISourceCompletionProvider interface {
	Native() uintptr
	toSourceCompletionProvider() *C.GtkSourceCompletionProvider

	GetName() string
	GetIcon() *gdk.Pixbuf
	GetIconName() string
	//GetGIcon() string {
	Populate(context *SourceCompletionContext)
	GetActivation() SourceCompletionActivation
	Match(context *SourceCompletionContext)
	GetInfoWidget(context *SourceCompletionProposal) (gtk.IWidget, error)
	UpdateInfo(*SourceCompletionProposal, *SourceCompletionInfo)
	GetStartIter(*SourceCompletionContext, *SourceCompletionProposal, *gtk.TextIter) bool
	ActivateProposal(*SourceCompletionProposal, *gtk.TextIter) bool
	GetInteractiveDelay() int
	GetPriority() int
}

// SourceCompletionProvider is a representation of GTK's GtkSourceCompletionProvider.
type SourceCompletionProvider struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceCompletionProvider.
func (v *SourceCompletionProvider) native() *C.GtkSourceCompletionProvider {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceCompletionProvider(p)
}

func marshalSourceCompletionProvider(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceCompletionProvider(obj), nil
}

func wrapSourceCompletionProvider(obj *glib.Object) *SourceCompletionProvider {
	return &SourceCompletionProvider{obj}
}

func (v *SourceCompletionProvider) toSourceCompletionProvider() *C.GtkSourceCompletionProvider {
	if v == nil {
		return nil
	}
	return v.native()
}

// GetName is a wrapper around gtk_source_completion_provider_get_name().
func (v *SourceCompletionProvider) GetName() string {
	c := C.gtk_source_completion_provider_get_name(v.native())
	return C.GoString((*C.char)(c))
}

// GetIcon is a wrapper around gtk_source_completion_provider_get_icon().
func (v *SourceCompletionProvider) GetIcon() *gdk.Pixbuf {
	c := C.gtk_source_completion_provider_get_icon(v.native())
	if c == nil {
		return nil
	}

	return &gdk.Pixbuf{glib.Take(unsafe.Pointer(c))}
}

// GetIconName is a wrapper around gtk_source_completion_provider_get_icon_name().
func (v *SourceCompletionProvider) GetIconName() string {
	c := C.gtk_source_completion_provider_get_icon_name(v.native())
	return C.GoString((*C.char)(c))
}

// TODO: need gotk3 gio.GIcon
// GetGIcon is a wrapper around gtk_source_completion_provider_get_gicon().
//func (v *SourceCompletionProvider) GetGIcon() string {
//}

// Populate is a wrapper around gtk_source_completion_provider_populate().
func (v *SourceCompletionProvider) Populate(context *SourceCompletionContext) {
	C.gtk_source_completion_provider_populate(v.native(), context.native())
}

// GetActivation is a wrapper around gtk_source_completion_provider_get_activation().
func (v *SourceCompletionProvider) GetActivation() SourceCompletionActivation {
	return SourceCompletionActivation(C.gtk_source_completion_provider_get_activation(v.native()))
}

// Match is a wrapper around gtk_source_completion_provider_match().
func (v *SourceCompletionProvider) Match(context *SourceCompletionContext) {
	C.gtk_source_completion_provider_match(v.native(), context.native())
}

// GetInfoWidget is a wrapper around gtk_source_completion_provider_get_info_widget().
func (v *SourceCompletionProvider) GetInfoWidget(proposal *SourceCompletionProposal) (gtk.IWidget, error) {
	w := C.gtk_source_completion_provider_get_info_widget(v.native(), proposal.native())
	if w == nil {
		return nil, nil
	}

	return castWidget(w)
}

// UpdateInfo is a wrapper around gtk_source_completion_provider_update_info().
func (v *SourceCompletionProvider) UpdateInfo(proposal *SourceCompletionProposal, info *SourceCompletionInfo) {
	C.gtk_source_completion_provider_update_info(v.native(), proposal.native(), info.native())
}

// GetStartIter is a wrapper around gtk_source_completion_provider_get_start_iter().
func (v *SourceCompletionProvider) GetStartIter(context *SourceCompletionContext, proposal *SourceCompletionProposal, iter *gtk.TextIter) bool {
	b := C.gtk_source_completion_provider_get_start_iter(v.native(), context.native(), proposal.native(), (*C.GtkTextIter)(unsafe.Pointer(iter)))
	return gobool(b)
}

// ActivateProposal is a wrapper around gtk_source_completion_provider_activate_proposal().
func (v *SourceCompletionProvider) ActivateProposal(proposal *SourceCompletionProposal, iter *gtk.TextIter) bool {
	b := C.gtk_source_completion_provider_activate_proposal(v.native(), proposal.native(), (*C.GtkTextIter)(unsafe.Pointer(iter)))
	return gobool(b)
}

// GetInteractiveDelay is a wrapper around gtk_source_completion_provider_get_interactive_delay().
func (v *SourceCompletionProvider) GetInteractiveDelay() int {
	return int(C.gtk_source_completion_provider_get_interactive_delay(v.native()))
}

// GetPriority is a wrapper around gtk_source_completion_provider_get_priority().
func (v *SourceCompletionProvider) GetPriority() int {
	return int(C.gtk_source_completion_provider_get_priority(v.native()))
}

/*
 * GtkSourceCompletionProposal
 */

// SourceCompletionProposal is a representation of GtkSourceCompletionProposal.
type SourceCompletionProposal struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceCompletionProposal.
func (v *SourceCompletionProposal) native() *C.GtkSourceCompletionProposal {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceCompletionProposal(p)
}

func marshalSourceCompletionProposal(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceCompletionProposal(obj), nil
}

func wrapSourceCompletionProposal(obj *glib.Object) *SourceCompletionProposal {
	return &SourceCompletionProposal{obj}
}

// GetLabel is a wrapper around gtk_source_completion_proposal_get_label().
func (v *SourceCompletionProposal) GetLabel() string {
	c := C.gtk_source_completion_proposal_get_label(v.native())
	return C.GoString((*C.char)(c))
}

// GetMarkup is a wrapper around gtk_source_completion_proposal_get_markup().
func (v *SourceCompletionProposal) GetMarkup() string {
	c := C.gtk_source_completion_proposal_get_markup(v.native())
	return C.GoString((*C.char)(c))
}

// GetText is a wrapper around gtk_source_completion_proposal_get_text().
func (v *SourceCompletionProposal) GetText() string {
	c := C.gtk_source_completion_proposal_get_text(v.native())
	return C.GoString((*C.char)(c))
}

// GetIcon is a wrapper around gtk_source_completion_proposal_get_icon().
func (v *SourceCompletionProposal) GetIcon() *gdk.Pixbuf {
	c := C.gtk_source_completion_proposal_get_icon(v.native())
	if c == nil {
		return nil
	}

	return &gdk.Pixbuf{glib.Take(unsafe.Pointer(c))}
}

// GetIconName is a wrapper around gtk_source_completion_proposal_get_icon_name().
func (v *SourceCompletionProposal) GetIconName() string {
	c := C.gtk_source_completion_proposal_get_icon_name(v.native())
	return C.GoString((*C.char)(c))
}

// TODO: need gotk3 gio.GIcon
// GetGIcon is a wrapper around gtk_source_completion_proposal_get_gicon().
//func (v *SourceCompletionProposal) GetGIcon() string {
//}

// GetInfo is a wrapper around gtk_source_completion_proposal_get_info().
func (v *SourceCompletionProposal) GetInfo() *SourceCompletionInfo {
	c := C.gtk_source_completion_proposal_get_info(v.native())
	return wrapSourceCompletionInfo(glib.Take(unsafe.Pointer(c)))
}

// Changed is a wrapper around gtk_source_completion_proposal_changed().
func (v *SourceCompletionProposal) Changed() {
	C.gtk_source_completion_proposal_changed(v.native())
}

// Hash is a wrapper around gtk_source_completion_proposal_hash().
func (v *SourceCompletionProposal) Hash() uint {
	return uint(C.gtk_source_completion_proposal_hash(v.native()))
}

// Equal is a wrapper around gtk_source_completion_proposal_equal().
func (v *SourceCompletionProposal) Equal(other *SourceCompletionProposal) bool {
	return gobool(C.gtk_source_completion_proposal_equal(v.native(), other.native()))
}

/*
 * GtkSourceGutter
 */

// SourceGutter is a representation of GtkSourceGutter.
type SourceGutter struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceGutter.
func (v *SourceGutter) native() *C.GtkSourceGutter {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceGutter(p)
}

func marshalSourceGutter(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceGutter(obj), nil
}

func wrapSourceGutter(obj *glib.Object) *SourceGutter {
	return &SourceGutter{obj}
}

/*
 * GtkSourceView
 */

// SourceView is a representation of GtkSourceView.
type SourceView struct {
	gtk.TextView
}

// SetHighlightCurrentLine is a wrapper around gtk_source_view_set_highlight_current_line().
func (v *SourceView) SetHighlightCurrentLine(highlight bool) {
	C.gtk_source_view_set_highlight_current_line(v.native(), gbool(highlight))
}

// SetShowLineNumbers is a wrapper around gtk_source_view_set_show_line_numbers().
func (v *SourceView) SetShowLineNumbers(show bool) {
	C.gtk_source_view_set_show_line_numbers(v.native(), gbool(show))
}

// SetShowRightMargin is a wrapper around gtk_source_view_get_show_right_margin().
func (v *SourceView) SetShowRightMargin() {
	C.gtk_source_view_get_show_right_margin(v.native())
}

// native returns a pointer to the underlying GtkSourceView.
func (v *SourceView) native() *C.GtkSourceView {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceView(p)
}

// native returns a pointer to the underlying GtkSourceView.
func (v *SourceView) asTextView() *C.GtkTextView {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkTextView(p)
}

func marshalSourceView(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceView(obj), nil
}

func wrapSourceView(obj *glib.Object) *SourceView {
	return &SourceView{gtk.TextView{gtk.Container{gtk.Widget{glib.InitiallyUnowned{obj}}}}}
}

// SourceViewNew is a wrapper around gtk_source_view_new().
func SourceViewNew() (*SourceView, error) {
	c := C.gtk_source_view_new()
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceView(glib.Take(unsafe.Pointer(c))), nil
}

func SourceViewNewWithBuffer(buffer *SourceBuffer) (*SourceView, error) {
	c := C.gtk_source_view_new_with_buffer(buffer.native())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceView(glib.Take(unsafe.Pointer(c))), nil
}

// GetBuffer is a wrapper around gtk_source_view_get_buffer().
func (v *SourceView) GetBuffer() (*SourceBuffer, error) {
	c := C.gtk_text_view_get_buffer(v.asTextView())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceBuffer(glib.Take(unsafe.Pointer(c))), nil
}

// GetCompletion is a wrapper around gtk_source_view_get_buffer().
func (v *SourceView) GetCompletion() (*SourceCompletion, error) {
	c := C.gtk_source_view_get_completion(v.native())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceCompletion(glib.Take(unsafe.Pointer(c))), nil
}

// GetGutter is a wrapper around gtk_source_view_get_gutter().
func (v *SourceView) GetGutter(wt gtk.TextWindowType) (*SourceGutter, error) {
	c := C.gtk_source_view_get_gutter(v.native(), C.GtkTextWindowType(wt))
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceGutter(glib.Take(unsafe.Pointer(c))), nil
}

/*
 * GtkSourceBuffer
 */

// SourceBuffer is a representation of GtkSourceBuffer.
type SourceBuffer struct {
	gtk.TextBuffer
}

// native returns a pointer to the underlying GtkSourceBuffer.
func (v *SourceBuffer) native() *C.GtkSourceBuffer {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceBuffer(p)
}

// native returns a pointer to the underlying GtkSourceBuffer.
func (v *SourceBuffer) asTextBuffer() *C.GtkTextBuffer {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkTextBuffer(p)
}

func marshalSourceBuffer(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceBuffer(obj), nil
}

func wrapSourceBuffer(obj *glib.Object) *SourceBuffer {
	return &SourceBuffer{gtk.TextBuffer{obj}}
}

// SourceBufferNew is a wrapper around gtk_text_buffer_new().
func SourceBufferNew() (*SourceBuffer, error) {
	c := C.gtk_text_buffer_new(nil)
	if c == nil {
		return nil, errNilPtr
	}

	e := wrapSourceBuffer(glib.Take(unsafe.Pointer(c)))
	return e, nil
}

// SourceBufferNewWithLanguage is a wrapper around gtk_source_buffer_new_with_language().
func SourceBufferNewWithLanguage(l *SourceLanguage) (*SourceBuffer, error) {
	c := C.gtk_source_buffer_new_with_language(l.native())
	if c == nil {
		return nil, errNilPtr
	}

	e := wrapSourceBuffer(glib.Take(unsafe.Pointer(c)))
	return e, nil
}

// SetText is a wrapper around gtk_text_buffer_set_text().
func (v *SourceBuffer) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_text_buffer_set_text(v.asTextBuffer(), (*C.gchar)(cstr),
		C.gint(len(text)))
}

// SetLanguage is a wrapper around gtk_source_buffer_set_language().
func (v *SourceBuffer) SetLanguage(l *SourceLanguage) {
	C.gtk_source_buffer_set_language(v.native(), l.native())
}

// BeginNotUndoableAction is a wrapper around gtk_source_buffer_begin_not_undoable_action().
func (v *SourceBuffer) BeginNotUndoableAction() {
	C.gtk_source_buffer_begin_not_undoable_action(v.native())
}

// EndNotUndoableAction is a wrapper around gtk_source_buffer_end_not_undoable_action().
func (v *SourceBuffer) EndNotUndoableAction() {
	C.gtk_source_buffer_end_not_undoable_action(v.native())
}

// GetMaxUndoLevels is a wrapper around gtk_source_buffer_get_max_undo_levels().
func (v *SourceBuffer) GetMaxUndoLevels() {
	C.gtk_source_buffer_get_max_undo_levels(v.native())
}

// SetMaxUndoLevels is a wrapper around gtk_source_buffer_set_max_undo_levels().
func (v *SourceBuffer) SetMaxUndoLevels(levels int) {
	C.gtk_source_buffer_set_max_undo_levels(v.native(), C.gint(levels))
}

// SetStyleScheme is a wrapper around gtk_source_buffer_set_style_scheme().
func (v *SourceBuffer) SetStyleScheme(scheme *SourceStyleScheme) {
	C.gtk_source_buffer_set_style_scheme(v.native(), scheme.native())
}

/*
 * GtkSourceLanguageManager
 */

// SourceLanguageManager is a representation of GtkSourceLanguageManager.
type SourceLanguageManager struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceLanguageManager.
func (v *SourceLanguageManager) native() *C.GtkSourceLanguageManager {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceLanguageManager(p)
}

func marshalSourceLanguageManager(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceLanguageManager(obj), nil
}

func wrapSourceLanguageManager(obj *glib.Object) *SourceLanguageManager {
	return &SourceLanguageManager{obj}
}

// SourceLanguageManagerNew is a wrapper around gtk_text_buffer_new().
func SourceLanguageManagerNew() (*SourceLanguageManager, error) {
	c := C.gtk_source_language_manager_new()
	if c == nil {
		return nil, errNilPtr
	}

	e := wrapSourceLanguageManager(glib.Take(unsafe.Pointer(c)))
	return e, nil
}

// SourceLanguageManagerGetDefault is a wrapper around gtk_source_language_manager_get_default().
func SourceLanguageManagerGetDefault() (*SourceLanguageManager, error) {
	c := C.gtk_source_language_manager_get_default()
	if c == nil {
		return nil, errNilPtr
	}

	e := wrapSourceLanguageManager(glib.Take(unsafe.Pointer(c)))
	return e, nil
}

// GetLanguage is a wrapper around gtk_source_language_manager_get_language().
func (v *SourceLanguageManager) GetLanguage(id string) (*SourceLanguage, error) {
	cstr := C.CString(id)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_source_language_manager_get_language(v.native(), (*C.gchar)(cstr))
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceLanguage(glib.Take(unsafe.Pointer(c))), nil
}

/*
 * GtkSourceLanguage
 */

// SourceLanguage is a representation of GtkSourceLanguage.
type SourceLanguage struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceLanguageManager.
func (v *SourceLanguage) native() *C.GtkSourceLanguage {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceLanguage(p)
}

func marshalSourceLanguage(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceLanguage(obj), nil
}

func wrapSourceLanguage(obj *glib.Object) *SourceLanguage {
	return &SourceLanguage{obj}
}

/*
 * GtkSourceStyle
 */

// SourceStyle is a representation of GtkSourceStyle.
type SourceStyle struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceStyle.
func (v *SourceStyle) native() *C.GtkSourceStyle {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceStyle(p)
}

func marshalSourceStyle(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyle(obj), nil
}

func wrapSourceStyle(obj *glib.Object) *SourceStyle {
	return &SourceStyle{obj}
}

// Copy is a wrapper around gtk_source_style_copy().
func (v *SourceStyle) Copy() (*SourceStyle, error) {
	c := C.gtk_source_style_copy(v.native())
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceStyle(glib.Take(unsafe.Pointer(c))), nil
}

// Apply is a wrapper around gtk_source_style_apply().
func (v *SourceStyle) Apply(tag *gtk.TextTag) {
	ctag := C.toGtkTextTag(unsafe.Pointer(tag.GObject))
	C.gtk_source_style_apply(v.native(), ctag)
}

/*
 * GtkSourceStyleScheme
 */

// SourceStyleScheme is a representation of GtkSourceStyleScheme.
type SourceStyleScheme struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceStyleScheme.
func (v *SourceStyleScheme) native() *C.GtkSourceStyleScheme {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceStyleScheme(p)
}

func marshalSourceStyleScheme(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyleScheme(obj), nil
}

func wrapSourceStyleScheme(obj *glib.Object) *SourceStyleScheme {
	return &SourceStyleScheme{obj}
}

// GetID is a wrapper around gtk_source_style_scheme_get_id().
func (v *SourceStyleScheme) GetID() (string, error) {
	c := C.gtk_source_style_scheme_get_id(v.native())
	if c == nil {
		return "", errNilPtr
	}
	gostr := goString(c)
	C.g_free(C.gpointer(c))
	return gostr, nil
}

// GetName is a wrapper around gtk_source_style_scheme_get_name().
func (v *SourceStyleScheme) GetName() (string, error) {
	c := C.gtk_source_style_scheme_get_name(v.native())
	if c == nil {
		return "", errNilPtr
	}
	gostr := goString(c)
	C.g_free(C.gpointer(c))
	return gostr, nil
}

// GetDescription is a wrapper around gtk_source_style_scheme_get_description().
func (v *SourceStyleScheme) GetDescription() (string, error) {
	c := C.gtk_source_style_scheme_get_description(v.native())
	if c == nil {
		return "", errNilPtr
	}
	gostr := goString(c)
	C.g_free(C.gpointer(c))
	return gostr, nil
}

// GetAuthors is a wrapper around gtk_source_style_scheme_get_authors().
func (v *SourceStyleScheme) GetAuthors() []string {
	var authors []string
	cauthors := C.gtk_source_style_scheme_get_authors(v.native())
	if cauthors == nil {
		return nil
	}
	for {
		if *cauthors == nil {
			break
		}
		authors = append(authors, C.GoString((*C.char)(*cauthors)))
		cauthors = C.next_gcharptr(cauthors)
	}
	return authors
}

// GetFileName is a wrapper around gtk_source_style_scheme_get_filename().
func (v *SourceStyleScheme) GetFileName() (string, error) {
	c := C.gtk_source_style_scheme_get_filename(v.native())
	if c == nil {
		return "", errNilPtr
	}
	gostr := goString(c)
	C.g_free(C.gpointer(c))
	return gostr, nil
}

// GetStyle is a wrapper around gtk_source_style_scheme_get_style().
func (v *SourceStyleScheme) GetStyle(id string) (*SourceStyle, error) {
	cstr1 := (*C.gchar)(C.CString(id))
	defer C.free(unsafe.Pointer(cstr1))

	c := C.gtk_source_style_scheme_get_style(v.native(), cstr1)
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceStyle(glib.Take(unsafe.Pointer(c))), nil
}

/*
 * GtkSourceStyleSchemeManager
 */

// SourceStyleSchemeManager is a representation of GtkSourceStyleSchemeManager.
type SourceStyleSchemeManager struct {
	*glib.Object
}

// native returns a pointer to the underlying GtkSourceStyleSchemeManager.
func (v *SourceStyleSchemeManager) native() *C.GtkSourceStyleSchemeManager {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceStyleSchemeManager(p)
}

func marshalSourceStyleSchemeManager(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyleSchemeManager(obj), nil
}

func wrapSourceStyleSchemeManager(obj *glib.Object) *SourceStyleSchemeManager {
	return &SourceStyleSchemeManager{obj}
}

// SourceStyleSchemeManagerNew is a wrapper around gtk_source_style_scheme_manager_new().
func SourceStyleSchemeManagerNew() (*SourceStyleSchemeManager, error) {
	c := C.gtk_source_style_scheme_manager_new()
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceStyleSchemeManager(glib.Take(unsafe.Pointer(c))), nil
}

// SourceStyleSchemeManagerGetDefault is a wrapper around gtk_source_style_scheme_manager_get_default().
func SourceStyleSchemeManagerGetDefault() (*SourceStyleSchemeManager, error) {
	c := C.gtk_source_style_scheme_manager_get_default()
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceStyleSchemeManager(glib.Take(unsafe.Pointer(c))), nil
}

// SetSearchPath is a wrapper around gtk_source_style_scheme_manager_set_search_path().
func (v *SourceStyleSchemeManager) SetSearchPath(paths []string) {
	cpaths := C.make_strings(C.int(len(paths) + 1))
	for i, path := range paths {
		cstr := C.CString(path)
		defer C.free(unsafe.Pointer(cstr))
		C.set_string(cpaths, C.int(i), (*C.gchar)(cstr))
	}

	C.set_string(cpaths, C.int(len(paths)), nil)
	C.gtk_source_style_scheme_manager_set_search_path(v.native(), cpaths)
	C.destroy_strings(cpaths)
}

// AppendSearchPath is a wrapper around gtk_source_style_scheme_manager_append_search_path().
func (v *SourceStyleSchemeManager) AppendSearchPath(path string) {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_source_style_scheme_manager_append_search_path(v.native(), (*C.gchar)(cstr))
}

// PrependSearchPath is a wrapper around gtk_source_style_scheme_manager_prepend_search_path().
func (v *SourceStyleSchemeManager) PrependSearchPath(path string) {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_source_style_scheme_manager_prepend_search_path(v.native(), (*C.gchar)(cstr))
}

// GetSearchPath is a wrapper around gtk_source_style_scheme_manager_get_search_path().
func (v *SourceStyleSchemeManager) GetSearchPath() []string {
	var paths []string
	cpaths := C.gtk_source_style_scheme_manager_get_search_path(v.native())
	if cpaths == nil {
		return nil
	}
	for {
		if *cpaths == nil {
			break
		}
		paths = append(paths, C.GoString((*C.char)(*cpaths)))
		cpaths = C.next_gcharptr(cpaths)
	}
	return paths
}

// GetSchemeIDs is a wrapper around gtk_source_style_scheme_manager_get_scheme_ids().
func (v *SourceStyleSchemeManager) GetSchemeIDs() []string {
	var ids []string
	cids := C.gtk_source_style_scheme_manager_get_scheme_ids(v.native())
	if cids == nil {
		return nil
	}
	for {
		if *cids == nil {
			break
		}
		ids = append(ids, C.GoString((*C.char)(*cids)))
		cids = C.next_gcharptr(cids)
	}
	return ids
}

// GetScheme is a wrapper around gtk_source_style_scheme_manager_get_scheme().
func (v *SourceStyleSchemeManager) GetScheme(id string) *SourceStyleScheme {
	cstr1 := (*C.gchar)(C.CString(id))
	defer C.free(unsafe.Pointer(cstr1))

	c := C.gtk_source_style_scheme_manager_get_scheme(v.native(), cstr1)
	if c == nil {
		return nil
	}
	return wrapSourceStyleScheme(glib.Take(unsafe.Pointer(c)))
}

/*
 * GtkSourceStyleSchemeChooser
 */

// ISourceStyleSchemeChooser is an interface type implemented by all structs
// embedding a GtkSourceStyleSchemeChooser.  It is meant to be used as an
// argument type for wrapper functions that wrap around a C GTK function taking
// a GtkSourceStyleSchemeChooser.
type ISourceStyleSchemeChooser interface {
	toSourceStyleSchemeChooser() *C.GtkSourceStyleSchemeChooser
}

/*
 * GtkSourceStyleSchemeChooser
 */

// SourceStyleSchemeChooser is a representation of GtkSourceView's
// GtkSourceStyleSchemeChooser GInterface.
type SourceStyleSchemeChooser struct {
	*glib.Object
}

// native returns a pointer to the underlying GObject as a GtkSourceStyleSchemeChooser.
func (v *SourceStyleSchemeChooser) native() *C.GtkSourceStyleSchemeChooser {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkSourceStyleSchemeChooser(p)
}

func (v *SourceStyleSchemeChooser) toSourceStyleSchemeChooser() *C.GtkSourceStyleSchemeChooser {
	if v == nil {
		return nil
	}
	return v.native()
}

func marshalSourceStyleSchemeChooser(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyleSchemeChooser(obj), nil
}

func wrapSourceStyleSchemeChooser(obj *glib.Object) *SourceStyleSchemeChooser {
	return &SourceStyleSchemeChooser{obj}
}

// GetScheme is a wrapper around gtk_source_style_scheme_chooser_get_style_scheme().
func (v *SourceStyleSchemeChooser) GetScheme() *SourceStyleScheme {
	c := C.gtk_source_style_scheme_chooser_get_style_scheme(v.native())
	if c == nil {
		return nil
	}
	return wrapSourceStyleScheme(glib.Take(unsafe.Pointer(c)))
}

// SetScheme is a wrapper around gtk_source_style_scheme_chooser_set_style_scheme().
func (v *SourceStyleSchemeChooser) SetScheme(scheme *SourceStyleScheme) {
	C.gtk_source_style_scheme_chooser_set_style_scheme(v.native(), scheme.native())
}

/*
 * GtkSourceStyleSchemeChooserButton
 */

// SourceStyleSchemeChooserButton is a representation of GtkSourceStyleSchemeChooserButton.
type SourceStyleSchemeChooserButton struct {
	gtk.Button

	SourceStyleSchemeChooser
}

// native returns a pointer to the underlying GtkSourceStyleSchemeChooserButton.
func (v *SourceStyleSchemeChooserButton) native() *C.GtkSourceStyleSchemeChooserButton {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toSourceStyleSchemeChooserButton(p)
}

func (v *SourceStyleSchemeChooserButton) toSourceStyleSchemeChooser() *C.GtkSourceStyleSchemeChooser {
	if v == nil {
		return nil
	}
	return C.toGtkSourceStyleSchemeChooser(unsafe.Pointer(v.GObject))
}

func marshalSourceStyleSchemeChooserButton(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyleSchemeChooserButton(obj), nil
}

func wrapSourceStyleSchemeChooserButton(obj *glib.Object) *SourceStyleSchemeChooserButton {
	actionable := &gtk.Actionable{obj}
	chooser := wrapSourceStyleSchemeChooser(obj)
	return &SourceStyleSchemeChooserButton{
		gtk.Button{
			gtk.Bin{
				gtk.Container{
					gtk.Widget{
						glib.InitiallyUnowned{obj},
					},
				},
			},
			actionable,
		},
		*chooser,
	}
}

/*
 * GtkSourceStyleSchemeChooserWidget
 */

// SourceStyleSchemeChooserWidget is a representation of GtkSourceStyleSchemeChooserWidget.
type SourceStyleSchemeChooserWidget struct {
	gtk.Bin

	SourceStyleSchemeChooser
}

// native returns a pointer to the underlying GtkSourceStyleSchemeChooserWidget.
func (v *SourceStyleSchemeChooserWidget) native() *C.GtkSourceStyleSchemeChooserWidget {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toSourceStyleSchemeChooserWidget(p)
}

func (v *SourceStyleSchemeChooserWidget) toSourceStyleSchemeChooser() *C.GtkSourceStyleSchemeChooser {
	if v == nil {
		return nil
	}
	return C.toGtkSourceStyleSchemeChooser(unsafe.Pointer(v.GObject))
}

func marshalSourceStyleSchemeChooserWidget(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapSourceStyleSchemeChooserWidget(obj), nil
}

func wrapSourceStyleSchemeChooserWidget(obj *glib.Object) *SourceStyleSchemeChooserWidget {
	chooser := wrapSourceStyleSchemeChooser(obj)
	return &SourceStyleSchemeChooserWidget{
		gtk.Bin{
			gtk.Container{
				gtk.Widget{
					glib.InitiallyUnowned{obj},
				},
			},
		},
		*chooser,
	}
}

// SourceStyleSchemeChooserWidgetNew is a wrapper around gtk_source_style_scheme_chooser_widget_new().
func SourceStyleSchemeChooserWidgetNew() (*SourceStyleSchemeChooserWidget, error) {
	c := C.gtk_source_style_scheme_chooser_widget_new()
	if c == nil {
		return nil, errNilPtr
	}
	return wrapSourceStyleSchemeChooserWidget(glib.Take(unsafe.Pointer(c))), nil
}
