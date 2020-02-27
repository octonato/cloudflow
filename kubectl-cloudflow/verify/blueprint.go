package verify

import (
	"fmt"
	"github.com/go-akka/configuration"
	"github.com/lightbend/cloudflow/kubectl-cloudflow/domain"
	"math/big"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// BlueprintProblem - generic interface for all blueprint related problems
type BlueprintProblem interface {
	ToMessage() string
}

type AmbiguousStreamletRef struct {
	BlueprintProblem
	streamletRef       string
	streamletClassName string
}

type BacktrackingVolumeMounthPath struct {
	BlueprintProblem
	className string
	name      string
	path      string
}

type InletProblem interface {
	BlueprintProblem
	InletPath() VerifiedPortPath
}

type PortPathError struct {
	BlueprintProblem
}

type InvalidStreamletName struct {
	BlueprintProblem
	streamletRef string
}

type InvalidStreamletClassName struct {
	BlueprintProblem
	streamletRef       string
	streamletClassName string
}

type StreamletDescriptorNotFound struct {
	BlueprintProblem
	streamletRef       string
	streamletClassName string
}

type DuplicateStreamletNamesFound struct {
	BlueprintProblem
	streamlets []StreamletRef
}

type InvalidConfigParameterKeyName struct {
	BlueprintProblem
	className string
	keyName   string
}

type InvalidValidationPatternConfigParameter struct {
	BlueprintProblem
	className         string
	keyName           string
	validationPattern string
}

type EmptyStreamlets struct {
	BlueprintProblem
}

type EmptyImages struct {
	BlueprintProblem
}

type EmptyStreamletDescriptors struct {
	BlueprintProblem
}

type DuplicateConfigParameterKeyFound struct {
	BlueprintProblem
	className string
	keyName   string
}

type DuplicateVolumeMountName struct {
	BlueprintProblem
	className string
	name      string
}

type DuplicateVolumeMountPath struct {
	BlueprintProblem
	className string
	path      string
}

type EmptyVolumeMountPath struct {
	BlueprintProblem
	className string
	name      string
}

type InvalidDefaultValueInConfigParameter struct {
	BlueprintProblem
	className    string
	keyName      string
	defaultValue string
}

type IllegalConnection struct {
	InletProblem
	outletPaths []VerifiedPortPath
	inletPath   VerifiedPortPath
}

type UnconnectedInlet struct {
	streamletRef string
	inlet        domain.InOutlet
}

type InvalidInletName struct {
	BlueprintProblem
	className string
	name      string
}

type InvalidOutletName struct {
	BlueprintProblem
	className string
	name      string
}

type IncompatibleSchema struct {
	InletProblem
	outletPortPath VerifiedPortPath
	inletPath      VerifiedPortPath
}

type InvalidPortPath struct {
	PortPathError
	path string
}

type InvalidVolumeMountName struct {
	BlueprintProblem
	className string
	name      string
}

type NonAbsoluteVolumeMountPath struct {
	BlueprintProblem
	className string
	name      string
	path      string
}

type PortPathNotFound struct {
	PortPathError
	path        string
	suggestions []VerifiedPortPath
}

type UnconnectedInlets struct {
	BlueprintProblem
	unconnectedInlets []UnconnectedInlet
}

type ImageInStreamletNotInImages struct {
	BlueprintProblem
	imageInStreamlet string
	streamlet        string
}

type StreamletNotInImageLabel struct {
	BlueprintProblem
	imageID   string
	streamlet string
}

func (b StreamletNotInImageLabel) ToMessage() string {
	return fmt.Sprintf("Streamlet %s not present in label of image %s", b.streamlet, b.imageID)
}

func (b ImageInStreamletNotInImages) ToMessage() string {
	return fmt.Sprintf("The image id %s referred to in streamlet %s does not appear in images section", b.imageInStreamlet, b.streamlet)
}

func (b AmbiguousStreamletRef) ToMessage() string {
	return fmt.Sprintf("ClassName matching %s is ambiguous for streamlet name %s.", b.streamletClassName, b.streamletRef)
}

func (b BacktrackingVolumeMounthPath) ToMessage() string {
	return fmt.Sprintf("`%s` contains a volume mount `%s` with an invalid path `$path`, backtracking in paths are not allowed.", b.className, b.name)
}

func (b InvalidStreamletName) ToMessage() string {
	return fmt.Sprintf("Invalid streamlet name %s. Names must consist of lower case alphanumeric characters and may contain '-' except for at the start or end.", b.streamletRef)
}

func (b InvalidStreamletClassName) ToMessage() string {
	return fmt.Sprintf("Class name %s for streamlet %s is invalid. Class names must be valid Java/Scala class names.", b.streamletClassName, b.streamletRef)
}

func (b StreamletDescriptorNotFound) ToMessage() string {
	return fmt.Sprintf("ClassName %s for %s cannot be found.", b.streamletClassName, b.streamletRef)
}

func (b DuplicateStreamletNamesFound) ToMessage() string {
	var duplicatesStreamlets = b.streamlets
	var duplicates string = ""
	for _, dup := range duplicatesStreamlets {
		if duplicates == "" {
			duplicates = fmt.Sprintf("(name: %s, className: %s)", dup.name, dup.className)
		} else {
			duplicates = duplicates + ", " + fmt.Sprintf("(name: %s, className: %s)", dup.name, dup.className)
		}
	}
	return fmt.Sprintf("Duplicate streamlet names detected: %s.", duplicates)
}

func (b InvalidConfigParameterKeyName) ToMessage() string {
	return fmt.Sprintf("`%s` contains a configuration parameter with invalid key name %s.", b.className, b.keyName)
}

func (b InvalidValidationPatternConfigParameter) ToMessage() string {
	return fmt.Sprintf("`%s` contains a configuration parameter `%s` with an invalid validation pattern `%s`.", b.className, b.keyName, b.validationPattern)
}

func (b EmptyStreamletDescriptors) ToMessage() string {
	return fmt.Sprintf("The streamlet descriptor list is empty.")
}

func (b EmptyStreamlets) ToMessage() string {
	return fmt.Sprintf("The streamlets list is empty.")
}

func (b EmptyImages) ToMessage() string {
	return fmt.Sprintf("The images section is empty.")
}

func (b DuplicateConfigParameterKeyFound) ToMessage() string {
	return fmt.Sprintf("`%s` contains a duplicate configuration parameter key, `%s` is used in more than one `ConfigParameter`", b.className, b.keyName)
}

func (b DuplicateVolumeMountName) ToMessage() string {
	return fmt.Sprintf("`%s` contains volume mounts with duplicate names (`%s`).", b.className, b.name)
}

func (b DuplicateVolumeMountPath) ToMessage() string {
	return fmt.Sprintf("`%s` contains volume mounts with duplicate paths (`%s`).", b.className, b.path)
}

func (b EmptyVolumeMountPath) ToMessage() string {
	return fmt.Sprintf("`%s` contains a volume mount `%s` with an empty path.", b.className, b.name)
}

func (b InvalidDefaultValueInConfigParameter) ToMessage() string {
	return fmt.Sprintf("`%s` contains a configuration parameter `%s` with an invalid default value, `%s` is invalid.", b.className, b.keyName, b.defaultValue)
}

func (b IllegalConnection) ToMessage() string {
	var outletPathsFormatted string = ""
	for _, outlet := range b.outletPaths {
		if outletPathsFormatted == "" {
			outletPathsFormatted = outlet.ToString()
		} else {
			outletPathsFormatted = outletPathsFormatted + "," + outlet.ToString()
		}
	}
	return fmt.Sprintf("Illegal connection, too many outlet paths (%s) are connected to inlet %s.", outletPathsFormatted, b.inletPath.ToString())
}

func (b IllegalConnection) InletPath() VerifiedPortPath {
	return b.inletPath
}

func (b IncompatibleSchema) ToMessage() string {
	return fmt.Sprintf("Outlet%s is not compatible with inlet %s.", b.outletPortPath.ToString(), b.inletPath.ToString())
}

func (b IncompatibleSchema) InletPath() VerifiedPortPath {
	return b.inletPath
}

func (b InvalidInletName) ToMessage() string {
	return fmt.Sprintf("Inlet `%s` in streamlet `%s` is invalid. Names must consist of lower case alphanumeric characters and may contain '-' except for at the start or end.",
		b.name, b.className)
}

func (b InvalidOutletName) ToMessage() string {
	return fmt.Sprintf("Outlet `%s` in streamlet `%s` is invalid. Names must consist of lower case alphanumeric characters and may contain '-' except for at the start or end.",
		b.name, b.className)
}

func (b InvalidPortPath) ToMessage() string {
	return fmt.Sprintf("'%s' is not a valid path to an outlet or an inlet.", b.path)
}

func (b InvalidVolumeMountName) ToMessage() string {
	return fmt.Sprintf("Volume mount `%s` in streamlet `%s` is invalid. Names must consist of lower case alphanumeric characters and may contain '-' except for at the start or end.",
		b.name, b.className)
}

func (b NonAbsoluteVolumeMountPath) ToMessage() string {
	return fmt.Sprintf("`%s` contains a volume mount `%s` with a non-absolute path (`%s`).", b.className, b.name, b.path)
}

func (b PortPathNotFound) ToMessage() string {
	var end = "."
	if b.suggestions != nil {
		//TODO: refactor this to a mkString function
		var suggestionsFormatted string = ""
		for _, suggestion := range b.suggestions {
			if suggestionsFormatted == "" {
				suggestionsFormatted = suggestion.ToString()
			} else {
				suggestionsFormatted = suggestionsFormatted + " or " + suggestion.ToString()
			}
		}
		end = fmt.Sprintf(", please try %s.", suggestionsFormatted)
	}
	return fmt.Sprintf("'%s' does not point to a known streamlet inlet or outlet%s", b.path, end)
}

func (b UnconnectedInlets) ToMessage() string {
	var listFormatted = ""
	for _, in := range b.unconnectedInlets {
		if listFormatted == "" {
			listFormatted = fmt.Sprintf("%s,%s", in.streamletRef, in.inlet.Name)
		} else {
			listFormatted = listFormatted + "," + fmt.Sprintf("%s,%s", in.streamletRef, in.inlet.Name)
		}
	}
	return fmt.Sprintf("Inlets (%s) are not connected.", listFormatted)
}

type Blueprint struct {
	images               map[string]domain.ImageReference
	streamlets           []StreamletRef
	connections          []StreamletConnection
	streamletDescriptorsPerImage map[string][]StreamletDescriptor
	globalProblems       []BlueprintProblem
	allProblems          []BlueprintProblem
}

// check for consistency between the image id mentioned in streamlets and the
// image ids present in images section of the blueprint
func checkStreamletImageConsistency(blueprint Blueprint) []BlueprintProblem {
	var problems []BlueprintProblem
	// check if the image ids in streamlets are present in the images section
	for _, streamlet := range blueprint.streamlets {
		img := *streamlet.imageId
		if _, ok := blueprint.images[img]; !ok {
			problems = append(problems, ImageInStreamletNotInImages{imageInStreamlet: img, streamlet: streamlet.name})
		}
	}
	return problems
}

// check if the streamlet is really in the label of the image that prefixes it in the blueprint
func checkStreamletImageLabelConsistency(blueprint Blueprint) []BlueprintProblem {
	var problems []BlueprintProblem
	for _, streamlet := range blueprint.streamlets {
		imageID := *streamlet.imageId
		if _, ok := blueprint.streamletDescriptorsPerImage[imageID]; !ok {
			problems = append(problems, StreamletNotInImageLabel{imageID: imageID, streamlet: streamlet.name})
		}
	}
	return problems
}

func (b Blueprint) verify() Blueprint {

	var illegalConnectionProblems, unconnectedInletProblems, portNameProblems, configParameterProblems, volumeMountProblems []BlueprintProblem

	var emptyImagesProblem *EmptyImages
	if len(b.images) == 0 {
		emptyImagesProblem = &EmptyImages{}
	}

	var emptyStreamletsProblem *EmptyStreamlets
	if len(b.streamlets) == 0 {
		emptyStreamletsProblem = &EmptyStreamlets{}
	}

	var imageInStreamletNotInImagesErrors []BlueprintProblem 
	imageInStreamletNotInImagesErrors = append(imageInStreamletNotInImagesErrors, checkStreamletImageConsistency(b) ...)

	var streamletNotInImageLabelErrors []BlueprintProblem 
	streamletNotInImageLabelErrors = append(streamletNotInImageLabelErrors, checkStreamletImageLabelConsistency(b) ...)

	// get all streamlet descriptors for all images
	var streamletDescriptors []StreamletDescriptor
	for _, desc := range b.streamletDescriptorsPerImage {
		streamletDescriptors = append(streamletDescriptors, desc ...)
	}

	var emptyStreamletDescriptorsProblem *EmptyStreamletDescriptors
	if len(streamletDescriptors) == 0 {
		emptyStreamletDescriptorsProblem = &EmptyStreamletDescriptors{}
	}

	var newStreamlets []StreamletRef
	var verifiedStreamlets []VerifiedStreamlet

	for _, ref := range b.streamlets {
		newStreamlets = append(newStreamlets, ref.verify(streamletDescriptors))
	}

	for _, streamlet := range newStreamlets {
		if streamlet.verified != nil {
			verifiedStreamlets = append(verifiedStreamlets, *streamlet.verified)
		}
	}

	var newConnections []StreamletConnection
	var verifiedConnections []VerifiedStreamletConnection

	for _, con := range b.connections {
		newConnections = append(newConnections, con.verify(verifiedStreamlets))
	}

	for _, verCon := range newConnections {
		if verCon.verified != nil {
			verifiedConnections = append(verifiedConnections, *verCon.verified)
		}
	}

	_, duplicatesProblem := b.verifyNoDuplicateStreamletNames(newStreamlets)

	portNameProblems = b.verifyPortNames(streamletDescriptors)
	configParameterProblems = b.verifyConfigParameters(streamletDescriptors)
	volumeMountProblems = b.verifyVolumeMounts(streamletDescriptors)

	verifiedConnections, conProblems := b.verifyUniqueInletConnections(verifiedConnections)

	for _, conProblem := range conProblems {
		illegalConnectionProblems = append(illegalConnectionProblems, conProblem)
	}

	var inletProblems []BlueprintProblem

	inletProblems = append(inletProblems, illegalConnectionProblems...)

	for _, newCon := range newConnections {
		for _, p := range newCon.problems {
			_, ok := p.(InletProblem)

			if ok {
				inletProblems = append(inletProblems, p)
			}
		}
	}
	var globalProblems []BlueprintProblem

	_, inletConProblems := b.verifyInletsConnected(verifiedStreamlets, verifiedConnections)

	for _, inletConProblem := range inletConProblems {
		var filteredUnconnectedInlets []UnconnectedInlet
		filteredUnconnectedInlets = filterUnconnectedInlets(inletProblems, inletConProblem.unconnectedInlets)

		if len(filteredUnconnectedInlets) > 0 {
			unconnectedInletProblems = append(unconnectedInletProblems, UnconnectedInlets{unconnectedInlets: filteredUnconnectedInlets})
		}
	}

	if emptyStreamletsProblem != nil {
		globalProblems = append(globalProblems, *emptyStreamletsProblem)
	}

	if emptyImagesProblem != nil {
		globalProblems = append(globalProblems, *emptyImagesProblem)
	}

	if emptyStreamletDescriptorsProblem != nil {
		globalProblems = append(globalProblems, *emptyStreamletDescriptorsProblem)
	}

	if duplicatesProblem != nil {
		globalProblems = append(globalProblems, duplicatesProblem)
	}

	if len(imageInStreamletNotInImagesErrors) > 0 {
		globalProblems = append(globalProblems, imageInStreamletNotInImagesErrors ...)
	}

	if len(streamletNotInImageLabelErrors) > 0 {
		globalProblems = append(globalProblems, streamletNotInImageLabelErrors ...)
	}

	var problems = [][]BlueprintProblem{illegalConnectionProblems, unconnectedInletProblems, portNameProblems, configParameterProblems, volumeMountProblems}
	for i := range problems {
		globalProblems = append(globalProblems, problems[i]...)
	}

	return Blueprint{
		streamlets: newStreamlets, 
		connections: newConnections, 
		streamletDescriptorsPerImage: 
		b.streamletDescriptorsPerImage, 
		globalProblems: globalProblems,
	}
}

func filterUnconnectedInlets(inletProblems []BlueprintProblem, unconnectedInlets []UnconnectedInlet) []UnconnectedInlet {
	var res []UnconnectedInlet
	for _, unconnectedInlet := range unconnectedInlets {
		for _, p := range inletProblems {
			inletProblem, ok := p.(InletProblem)
			if ok {
				if !reflect.DeepEqual(inletProblem.InletPath(), VerifiedPortPath{streamletRef: unconnectedInlet.streamletRef, portName: &unconnectedInlet.inlet.Name}) {
					res = append(res, unconnectedInlet)
				}
			}
		}
	}
	return res
}

type GroupedConnections struct {
	vInlet VerifiedInlet
	vCons  []VerifiedStreamletConnection
}

func (b Blueprint) UpdateAllProblems() []BlueprintProblem {
	var streamletProblems []BlueprintProblem
	var connectionProblems []BlueprintProblem

	for _, streamlet := range b.streamlets {
		streamletProblems = append(streamletProblems, streamlet.problems...)
	}

	for _, connection := range b.connections {
		connectionProblems = append(connectionProblems, connection.problems...)
	}

	var problems = [][]BlueprintProblem{b.globalProblems, streamletProblems, connectionProblems}
	var res []BlueprintProblem
	for i := range problems {
		res = append(res, problems[i]...)
	}
	b.globalProblems = res
	return b.globalProblems
}

func (b Blueprint) verifyUniqueInletConnections(verifiedStreamletConnections []VerifiedStreamletConnection) ([]VerifiedStreamletConnection, []IllegalConnection) {
	groupedConnections := make(map[string]GroupedConnections)
	for i := range verifiedStreamletConnections {
		// cannot use a VerifiedInlet a a map key here
		hash := GetSHA256Hash(verifiedStreamletConnections[i].verifiedInlet)
		key := hash
		if val, ok := groupedConnections[key]; ok {
			val.vCons = append(groupedConnections[key].vCons, verifiedStreamletConnections[i])
		} else {
			values := []VerifiedStreamletConnection{}
			values = append(values, verifiedStreamletConnections[i])
			groupedConnections[key] = GroupedConnections{vInlet: verifiedStreamletConnections[i].verifiedInlet, vCons: values}
		}
	}
	var illegalConnectionProblems []IllegalConnection
	for _, gCon := range groupedConnections {
		if len(gCon.vCons) > 1 {
			var mapPortpaths []VerifiedPortPath
			for _, vOutlet := range gCon.vCons {
				mapPortpaths = append(mapPortpaths, vOutlet.verifiedOutlet.portPath())
			}
			illegalConnectionProblems = append(illegalConnectionProblems, IllegalConnection{
				outletPaths: mapPortpaths,
				inletPath:   gCon.vInlet.portPath(),
			})
		}
	}

	if len(illegalConnectionProblems) != 0 {
		return nil, illegalConnectionProblems
	} else {
		return verifiedStreamletConnections, nil
	}
}

func verifiedConnectionsExists(verifiedStreamletConnections []VerifiedStreamletConnection, inlet domain.InOutlet, streamlet VerifiedStreamlet) bool {
	for _, con := range verifiedStreamletConnections {
		if reflect.DeepEqual(con.verifiedInlet.streamlet, streamlet) && con.verifiedInlet.portName == inlet.Name {
			return true
		}
	}
	return false
}

func (b Blueprint) verifyInletsConnected(verifiedStreamlets []VerifiedStreamlet, verifiedStreamletConnections []VerifiedStreamletConnection) ([]VerifiedStreamlet, []UnconnectedInlets) {
	var unconnectedPortProblems []UnconnectedInlets

	for _, vStreamlet := range verifiedStreamlets {
		var unconnectedInlets []UnconnectedInlet

		for _, inlet := range vStreamlet.descriptor.Inlets {
			if !verifiedConnectionsExists(verifiedStreamletConnections, inlet, vStreamlet) {
				unconnectedInlets = append(unconnectedInlets, UnconnectedInlet{vStreamlet.name, inlet})
			}
		}

		if len(unconnectedInlets) != 0 {
			unconnectedPortProblems = append(unconnectedPortProblems, UnconnectedInlets{unconnectedInlets: unconnectedInlets})
		}
	}

	if len(unconnectedPortProblems) != 0 {
		return verifiedStreamlets, nil
	} else {
		return nil, unconnectedPortProblems
	}
}

func (b Blueprint) validate() (*Blueprint, []BlueprintProblem) {
	if len(b.allProblems) == 0 {
		return &b, nil
	} else {
		return nil, b.allProblems
	}
}

func (b Blueprint) verifyNoDuplicateStreamletNames(streamlets []StreamletRef) ([]StreamletRef, *DuplicateStreamletNamesFound) {
	groupedStreamlets := make(map[string][]StreamletRef)

	for i := range streamlets {
		sName := strings.TrimSpace(streamlets[i].name)
		if _, ok := groupedStreamlets[sName]; ok {
			groupedStreamlets[sName] = append(groupedStreamlets[sName], streamlets[i])
		} else {
			values := []StreamletRef{}
			values = append(values, streamlets[i])
			groupedStreamlets[sName] = values
		}
	}

	var duplicateStreamlets []StreamletRef
	for _, v := range groupedStreamlets {
		if len(v) > 1 {
			for _, ref := range v {
				duplicateStreamlets = append(duplicateStreamlets, ref)
			}

		}
	}

	if len(duplicateStreamlets) == 0 {
		return streamlets, nil
	} else {
		return nil, &DuplicateStreamletNamesFound{streamlets: duplicateStreamlets}
	}
}

func (b Blueprint) verifyPortNames(streamletDescriptors []StreamletDescriptor) []BlueprintProblem {
	var inletProblems []BlueprintProblem
	var outletProblems []BlueprintProblem

	for _, desc := range streamletDescriptors {
		for _, inlet := range desc.Inlets {
			if !IsDnsLabelCompatible(inlet.Name) {
				inletProblems = append(inletProblems, InvalidInletName{className: desc.ClassName, name: inlet.Name})
			}
		}

		for _, outlet := range desc.Outlets {
			if !IsDnsLabelCompatible(outlet.Name) {
				outletProblems = append(outletProblems, InvalidInletName{className: desc.ClassName, name: outlet.Name})
			}
		}
	}
	return append(inletProblems, outletProblems...)
}

func (b Blueprint) verifyVolumeMounts(streamletDescriptors []StreamletDescriptor) []BlueprintProblem {
	separator := string(filepath.Separator)
	var invalidPaths, invalidNames, duplicateNames, duplicatePaths []BlueprintProblem
	var names, paths []string
	for _, desc := range streamletDescriptors {
		for _, vMount := range desc.VolumeMounts {
			names = append(names, vMount.Name)
			paths = append(paths, vMount.Path)
			for _, path := range strings.Split(vMount.Path, separator) {
				if path == ".." {
					invalidPaths = append(invalidPaths, BacktrackingVolumeMounthPath{className: desc.ClassName, name: vMount.Name, path: vMount.Path})
					break
				}
			}

			if len(vMount.Path) == 0 {
				invalidPaths = append(invalidPaths, EmptyVolumeMountPath{className: desc.ClassName, name: vMount.Name})
			}

			if !filepath.IsAbs(vMount.Path) {
				invalidPaths = append(invalidPaths, NonAbsoluteVolumeMountPath{className: desc.ClassName, name: vMount.Name, path: vMount.Path})
			}

			if IsDnsLabelCompatible(vMount.Name) {
				if len(vMount.Name) > DNS1123LabelMaxLength {
					invalidNames = append(invalidNames, InvalidVolumeMountName{className: desc.ClassName, name: vMount.Name})
				}

			} else {
				invalidNames = append(invalidNames, InvalidVolumeMountName{className: desc.ClassName, name: vMount.Name})
			}

		}

		dupsNames := Distinct(Diff(names, Distinct(names)))

		for _, dupName := range dupsNames {
			duplicateNames = append(duplicateNames, DuplicateVolumeMountName{className: desc.ClassName, name: dupName})
		}

		dupsPaths := Distinct(Diff(paths, Distinct(paths)))

		for _, dupPath := range dupsPaths {
			duplicatePaths = append(duplicatePaths, DuplicateVolumeMountName{className: desc.ClassName, name: dupPath})
		}
	}

	var problems = [][]BlueprintProblem{invalidPaths, invalidNames, duplicateNames, duplicatePaths}
	var res []BlueprintProblem

	for i := range problems {
		res = append(res, problems[i]...)
	}

	return res
}

func getDurationFromConfig(config string) (duration time.Duration, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Parsing duration error: %v", r)
			}
		}
	}()
	return configuration.ParseString(fmt.Sprintf("value=%s", config)).GetTimeDuration("value", time.Nanosecond), nil
}

func getMemorySizeFromConfig(config string) (size *big.Int, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Parsing memory size error: %v", r)
			}
		}
	}()
	return configuration.ParseString(fmt.Sprintf("value=%s", config)).GetByteSize("value"), nil
}

func (b Blueprint) verifyConfigParameters(streamletDescriptors []StreamletDescriptor) []BlueprintProblem {
	var invalidConfigParametersKeyProblems, invalidDefaultValueOrPatternProblems, duplicateConfigParametersKeysFound []BlueprintProblem
	var keys []string
	for _, desc := range streamletDescriptors {
		for _, configParam := range desc.ConfigParameters {
			keys = append(keys, configParam.Key)

			if !CheckFullPatternMatch(configParam.Key, ConfigParameterKeyPattern) {
				invalidConfigParametersKeyProblems = append(invalidConfigParametersKeyProblems, InvalidConfigParameterKeyName{className: desc.ClassName, keyName: configParam.Key})
			}

			switch configParam.Type {
			case "string":
				reg, err := regexp.Compile(configParam.Pattern)
				if err != nil {
					invalidDefaultValueOrPatternProblems = append(invalidDefaultValueOrPatternProblems, InvalidValidationPatternConfigParameter{className: desc.ClassName, keyName: configParam.Key, validationPattern: configParam.Pattern})
				} else {
					if !CheckFullPatternMatch(configParam.DefaultValue, reg) {
						invalidDefaultValueOrPatternProblems = append(invalidDefaultValueOrPatternProblems, InvalidDefaultValueInConfigParameter{className: desc.ClassName, keyName: configParam.Key, defaultValue: configParam.DefaultValue})
					}
				}
			case "duration":
				_, err := getDurationFromConfig(configParam.DefaultValue)

				if err != nil {
					invalidDefaultValueOrPatternProblems = append(invalidDefaultValueOrPatternProblems, InvalidDefaultValueInConfigParameter{className: desc.ClassName, keyName: configParam.Key, defaultValue: configParam.DefaultValue})
				}
			case "memorysize":
				_, err := getMemorySizeFromConfig(configParam.DefaultValue)

				if err != nil {
					invalidDefaultValueOrPatternProblems = append(invalidDefaultValueOrPatternProblems, InvalidDefaultValueInConfigParameter{className: desc.ClassName, keyName: configParam.Key, defaultValue: configParam.DefaultValue})
				}
			}

		}
		dupKeys := Distinct(Diff(keys, Distinct(keys)))
		for _, dupKey := range dupKeys {
			duplicateConfigParametersKeysFound = append(duplicateConfigParametersKeysFound, DuplicateConfigParameterKeyFound{className: desc.ClassName, keyName: dupKey})
		}
	}

	var problems = [][]BlueprintProblem{invalidConfigParametersKeyProblems, invalidDefaultValueOrPatternProblems, duplicateConfigParametersKeysFound}
	var res []BlueprintProblem

	for i := range problems {
		res = append(res, problems[i]...)
	}

	return res
}
